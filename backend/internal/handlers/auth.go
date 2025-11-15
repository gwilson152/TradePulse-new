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
	"github.com/tradepulse/api/internal/models"
)

type RequestMagicLinkInput struct {
	Email string `json:"email"`
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
				logger.Error("Failed to create user", "error", err, "email", input.Email)
				writeError(w, http.StatusInternalServerError, "USER_CREATE_ERROR", "Failed to create user")
				return
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
