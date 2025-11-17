package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RequestMagicLinkInput struct {
	Email string `json:"email"`
}

type SignupWithPlanInput struct {
	Email    string `json:"email"`
	PlanType string `json:"plan_type"`
}

type VerifyMagicLinkResponse struct {
	JWT  string      `json:"jwt"`
	User models.User `json:"user"`
}

// RequestMagicLink handles magic link generation and email sending
func RequestMagicLink(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input RequestMagicLinkInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
			return
		}

		if input.Email == "" {
			writeError(w, http.StatusBadRequest, "INVALID_EMAIL", "Email is required")
			return
		}

		// Check if user exists, create if not
		user, err := db.GetUserByEmail(r.Context(), input.Email)
		if err != nil {
			// User doesn't exist, create new user
			now := time.Now()
			user = &models.User{
				ID:        uuid.New(),
				Email:     input.Email,
				CreatedAt: now,
				LastLogin: &now,
			}

			if err := db.CreateUser(r.Context(), user); err != nil {
				// Check if error is due to duplicate email (race condition)
				if err.Error() == "failed to create user: pq: duplicate key value violates unique constraint \"users_email_key\"" {
					// User was just created, fetch them
					user, err = db.GetUserByEmail(r.Context(), input.Email)
					if err != nil {
						logger.Error("Failed to get user after duplicate key error", "error", err, "email", input.Email)
						writeError(w, http.StatusInternalServerError, "USER_FETCH_ERROR", "Failed to get user")
						return
					}
				} else {
					logger.Error("Failed to create user", "error", err, "email", input.Email)
					writeError(w, http.StatusInternalServerError, "USER_CREATE_ERROR", "Failed to create user")
					return
				}
			}
		}

		// Generate magic link token
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			logger.Error("Failed to generate token", "error", err)
			writeError(w, http.StatusInternalServerError, "TOKEN_GENERATE_ERROR", "Failed to generate token")
			return
		}
		token := hex.EncodeToString(tokenBytes)

		// Store token with 15-minute expiry
		expiresAt := time.Now().Add(15 * time.Minute)
		if err := db.StoreMagicLinkToken(r.Context(), user.ID, token, expiresAt); err != nil {
			logger.Error("Failed to store token", "error", err, "userId", user.ID)
			writeError(w, http.StatusInternalServerError, "TOKEN_STORE_ERROR", "Failed to store token")
			return
		}

		// TODO: Send email with magic link
		// For now, just log it
		logger.Info("Magic link generated",
			"email", input.Email,
			"token", token,
			"link", "https://tradepulse.drivenw.com/auth/verify?token="+token)

		writeSuccess(w, http.StatusOK, map[string]string{
			"message": "Magic link sent to your email",
		})
	}
}

// SignupWithPlan handles user signup with plan selection (Beta - all plans free)
func SignupWithPlan(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input SignupWithPlanInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
			return
		}

		if input.Email == "" {
			writeError(w, http.StatusBadRequest, "INVALID_EMAIL", "Email is required")
			return
		}

		if input.PlanType == "" {
			input.PlanType = "starter" // Default plan
		}

		// Validate plan type
		validPlans := map[string]bool{
			"starter": true,
			"pro":     true,
			"premium": true,
		}
		if !validPlans[input.PlanType] {
			writeError(w, http.StatusBadRequest, "INVALID_PLAN", "Invalid plan type")
			return
		}

		// Check if user already exists
		existingUser, err := db.GetUserByEmail(r.Context(), input.Email)
		if err == nil && existingUser != nil {
			writeError(w, http.StatusConflict, "USER_EXISTS", "User already exists. Please sign in instead.")
			return
		}

		// Create new user with selected plan
		now := time.Now()
		user := &models.User{
			ID:             uuid.New(),
			Email:          input.Email,
			PlanType:       input.PlanType,
			PlanStatus:     "beta_free",
			PlanSelectedAt: &now,
			CreatedAt:      now,
			LastLogin:      &now,
		}

		if err := db.CreateUser(r.Context(), user); err != nil {
			logger.Error("Failed to create user", "error", err, "email", input.Email)
			writeError(w, http.StatusInternalServerError, "USER_CREATE_ERROR", "Failed to create user")
			return
		}

		// Generate magic link token
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			logger.Error("Failed to generate token", "error", err)
			writeError(w, http.StatusInternalServerError, "TOKEN_GENERATE_ERROR", "Failed to generate token")
			return
		}
		token := hex.EncodeToString(tokenBytes)

		// Store token with 15-minute expiry
		expiresAt := time.Now().Add(15 * time.Minute)
		if err := db.StoreMagicLinkToken(r.Context(), user.ID, token, expiresAt); err != nil {
			logger.Error("Failed to store token", "error", err, "userId", user.ID)
			writeError(w, http.StatusInternalServerError, "TOKEN_STORE_ERROR", "Failed to store token")
			return
		}

		// TODO: Send email with magic link
		// For now, just log it
		logger.Info("Signup with plan - Magic link generated",
			"email", input.Email,
			"plan", input.PlanType,
			"token", token,
			"link", "https://tradepulse.drivenw.com/auth/verify?token="+token)

		writeSuccess(w, http.StatusOK, map[string]string{
			"message": "Account created! Magic link sent to your email",
		})
	}
}

// VerifyMagicLink verifies the token and returns a JWT
func VerifyMagicLink(db *database.DB, logger *slog.Logger, jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			writeError(w, http.StatusBadRequest, "INVALID_TOKEN", "Token is required")
			return
		}

		// Verify and consume token
		userID, err := db.VerifyMagicLinkToken(r.Context(), token)
		if err != nil {
			logger.Error("Failed to verify token", "error", err, "token", token)
			writeError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		// Get user
		user, err := db.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error("Failed to get user", "error", err, "userId", userID)
			writeError(w, http.StatusInternalServerError, "USER_FETCH_ERROR", "Failed to get user")
			return
		}

		// Update last login
		if err := db.UpdateUserLastLogin(r.Context(), userID); err != nil {
			logger.Warn("Failed to update last login", "error", err, "userId", userID)
		}

		// Generate JWT
		expiryDuration, err := time.ParseDuration(jwtExpiry)
		if err != nil {
			expiryDuration = 24 * time.Hour
		}

		claims := jwt.MapClaims{
			"user_id": userID.String(),
			"email":   user.Email,
			"exp":     time.Now().Add(expiryDuration).Unix(),
			"iat":     time.Now().Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := jwtToken.SignedString([]byte(jwtSecret))
		if err != nil {
			logger.Error("Failed to sign JWT", "error", err)
			writeError(w, http.StatusInternalServerError, "JWT_SIGN_ERROR", "Failed to generate JWT")
			return
		}

		logger.Info("User authenticated", "userId", userID, "email", user.Email)

		writeSuccess(w, http.StatusOK, VerifyMagicLinkResponse{
			JWT:  signedToken,
			User: *user,
		})
	}
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		// Get user from database
		user, err := db.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error("Failed to get user", "error", err, "userId", userID)
			writeError(w, http.StatusInternalServerError, "USER_FETCH_ERROR", "Failed to get user")
			return
		}

		writeSuccess(w, http.StatusOK, user)
	}
}

type SetPasswordInput struct {
	Password string `json:"password"`
}

type LoginWithPasswordInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SetPassword allows users to set or update their password
func SetPassword(db *database.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := middleware.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			return
		}

		var input SetPasswordInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
			return
		}

		// Validate password
		if len(input.Password) < 8 {
			writeError(w, http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 8 characters")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Failed to hash password", "error", err)
			writeError(w, http.StatusInternalServerError, "HASH_ERROR", "Failed to process password")
			return
		}

		// Store password hash
		if err := db.SetUserPassword(r.Context(), userID, string(hashedPassword)); err != nil {
			logger.Error("Failed to set password", "error", err, "userId", userID)
			writeError(w, http.StatusInternalServerError, "PASSWORD_SET_ERROR", "Failed to set password")
			return
		}

		logger.Info("Password set successfully", "userId", userID)

		writeSuccess(w, http.StatusOK, map[string]string{
			"message": "Password set successfully",
		})
	}
}

// LoginWithPassword handles email/password authentication
func LoginWithPassword(db *database.DB, logger *slog.Logger, jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input LoginWithPasswordInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
			return
		}

		if input.Email == "" || input.Password == "" {
			writeError(w, http.StatusBadRequest, "INVALID_CREDENTIALS", "Email and password are required")
			return
		}

		// Get user by email
		user, err := db.GetUserByEmail(r.Context(), input.Email)
		if err != nil {
			logger.Warn("Login attempt for non-existent user", "email", input.Email)
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}

		// Check if user has a password set
		if !user.HasPassword || user.PasswordHash == "" {
			logger.Warn("Login attempt for user without password", "email", input.Email, "userId", user.ID)
			writeError(w, http.StatusUnauthorized, "NO_PASSWORD_SET", "No password set for this account. Please use magic link login.")
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			logger.Warn("Invalid password attempt", "email", input.Email, "userId", user.ID)
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}

		// Update last login
		if err := db.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			logger.Warn("Failed to update last login", "error", err, "userId", user.ID)
		}

		// Generate JWT
		expiryDuration, err := time.ParseDuration(jwtExpiry)
		if err != nil {
			expiryDuration = 24 * time.Hour
		}

		claims := jwt.MapClaims{
			"user_id": user.ID.String(),
			"email":   user.Email,
			"exp":     time.Now().Add(expiryDuration).Unix(),
			"iat":     time.Now().Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := jwtToken.SignedString([]byte(jwtSecret))
		if err != nil {
			logger.Error("Failed to sign JWT", "error", err)
			writeError(w, http.StatusInternalServerError, "JWT_SIGN_ERROR", "Failed to generate JWT")
			return
		}

		logger.Info("User authenticated with password", "userId", user.ID, "email", user.Email)

		writeSuccess(w, http.StatusOK, VerifyMagicLinkResponse{
			JWT:  signedToken,
			User: *user,
		})
	}
}
