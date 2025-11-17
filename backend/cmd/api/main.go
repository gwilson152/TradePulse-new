package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/handlers"
	appMiddleware "github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/notifications"
)

type application struct {
	db              *database.DB
	logger          *slog.Logger
	config          config
	notificationBus *notifications.Bus
}

type config struct {
	port           string
	environment    string
	allowedOrigins string
	jwtSecret      string
	jwtExpiry      string
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg := config{
		port:           getEnv("PORT", "9000"),
		environment:    getEnv("ENVIRONMENT", "development"),
		allowedOrigins: getEnv("ALLOWED_ORIGINS", "https://tradepulse.drivenw.com"),
		jwtSecret:      getEnv("JWT_SECRET", ""),
		jwtExpiry:      getEnv("JWT_EXPIRY", "24h"),
	}

	if cfg.jwtSecret == "" {
		logger.Error("JWT_SECRET environment variable is required")
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New(database.Config{
		Host:     getEnv("DB_HOST", "postgres1.drivenw.local"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "tradepulse"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "tradepulse"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Run migrations
	migrationsPath := getEnv("MIGRATIONS_PATH", "migrations")
	if err := db.RunMigrations(migrationsPath); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	logger.Info("Database migrations completed")

	// Initialize notification bus
	notificationBus := notifications.NewBus(logger)
	go notificationBus.Run()

	logger.Info("Notification bus started")

	// Initialize application
	app := &application{
		db:              db,
		logger:          logger,
		config:          cfg,
		notificationBus: notificationBus,
	}

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	allowedOrigins := []string{cfg.allowedOrigins}
	logger.Info("CORS configured", "allowedOrigins", allowedOrigins)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Initialize handlers
	tradesHandler := handlers.NewTradesHandler(app.db, app.notificationBus)
	tagsHandler := handlers.NewTagsHandler(app.db)
	csvImportHandler := handlers.NewCSVImportHandler(app.db, app.notificationBus)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes (no auth required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/request-magic-link", handlers.RequestMagicLink(app.db, app.logger))
			r.Post("/signup", handlers.SignupWithPlan(app.db, app.logger))
			r.Get("/verify", handlers.VerifyMagicLink(app.db, app.logger, cfg.jwtSecret, cfg.jwtExpiry))
			r.Post("/login", handlers.LoginWithPassword(app.db, app.logger, cfg.jwtSecret, cfg.jwtExpiry))
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Authenticate(cfg.jwtSecret))

			// Auth
			r.Get("/auth/me", handlers.GetCurrentUser(app.db, app.logger))
			r.Post("/auth/set-password", handlers.SetPassword(app.db, app.logger))
			r.Post("/auth/logout", handlers.Logout(app.logger))
			r.Post("/auth/refresh", handlers.RefreshToken(app.logger, cfg.jwtSecret, cfg.jwtExpiry))

			// Trades
			r.Get("/trades", tradesHandler.ListTrades)
			r.Post("/trades", tradesHandler.CreateTrade)
			r.Get("/trades/{id}", tradesHandler.GetTrade)
			r.Put("/trades/{id}", tradesHandler.UpdateTrade)
			r.Delete("/trades/{id}", tradesHandler.DeleteTrade)
			r.Post("/trades/import-csv", csvImportHandler.ImportCSV)

			// Trade tags
			r.Post("/trades/{id}/tags", tradesHandler.AddTagToTrade)
			r.Delete("/trades/{tradeId}/tags/{tagId}", tradesHandler.RemoveTagFromTrade)

			// Journal
			r.Get("/journal", handlers.ListJournalEntries(app.db, app.logger))
			r.Post("/journal", handlers.CreateJournalEntry(app.db, app.logger))
			r.Get("/journal/{id}", handlers.GetJournalEntry(app.db, app.logger))
			r.Put("/journal/{id}", handlers.UpdateJournalEntry(app.db, app.logger))
			r.Delete("/journal/{id}", handlers.DeleteJournalEntry(app.db, app.logger))

			// Journal entries by trade
			r.Get("/trades/{tradeId}/journal", handlers.GetJournalEntriesByTradeID(app.db, app.logger))

			// Attachments
			r.Post("/journal/{id}/attachments", handlers.UploadAttachment(app.db, app.logger))
			r.Get("/attachments/{id}", handlers.GetAttachment(app.db, app.logger))
			r.Delete("/attachments/{id}", handlers.DeleteAttachment(app.db, app.logger))

			// Tags
			r.Get("/tags", tagsHandler.ListTags)
			r.Post("/tags", tagsHandler.CreateTag)
			r.Get("/tags/{id}", tagsHandler.GetTag)
			r.Put("/tags/{id}", tagsHandler.UpdateTag)
			r.Delete("/tags/{id}", tagsHandler.DeleteTag)

			// Rule Sets
			r.Get("/rulesets", handlers.ListRuleSets(app.db, app.logger))
			r.Post("/rulesets", handlers.CreateRuleSet(app.db, app.logger))
			r.Get("/rulesets/{id}", handlers.GetRuleSet(app.db, app.logger))
			r.Put("/rulesets/{id}", handlers.UpdateRuleSet(app.db, app.logger))
			r.Delete("/rulesets/{id}", handlers.DeleteRuleSet(app.db, app.logger))

			// Rules within a ruleset
			r.Post("/rulesets/{ruleSetId}/rules", handlers.CreateRule(app.db, app.logger))
			r.Put("/rulesets/{ruleSetId}/rules/{ruleId}", handlers.UpdateRule(app.db, app.logger))
			r.Delete("/rulesets/{ruleSetId}/rules/{ruleId}", handlers.DeleteRule(app.db, app.logger))

			// Metrics
			r.Get("/metrics/summary", handlers.GetSummaryMetrics(app.db, app.logger))
			r.Get("/metrics/by-symbol", handlers.GetMetricsBySymbol(app.db, app.logger))
			r.Get("/metrics/daily", handlers.GetDailyPerformance(app.db, app.logger))

			// WebSocket notifications
			r.Get("/ws", handlers.HandleWebSocket(app.notificationBus, app.logger))

			// Notification stats
			r.Get("/notifications/stats", handlers.HandleNotificationStats(app.notificationBus, app.logger))

			// Integrations
			r.Post("/integrations/propreports/fetch", handlers.FetchPropReportsTrades(app.logger))
		})
	})

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.port,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.port, "environment", cfg.environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("TradePulse API server started successfully", "url", fmt.Sprintf("https://api.tradepulse.drivenw.com:%s", cfg.port))

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited properly")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
