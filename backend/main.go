package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anish-chanda/ferna/internal/auth"
	"github.com/anish-chanda/ferna/internal/db"
	"github.com/anish-chanda/ferna/internal/handlers"
	"github.com/anish-chanda/ferna/internal/logger"
	"github.com/anish-chanda/ferna/migrations"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	authpkg "github.com/go-pkgz/auth/v2"
)

// App holds the application dependencies
type App struct {
	config *Config
	logger *logger.ServiceLogger
	db     *db.PostgresDB
	server *http.Server
	auth   *authpkg.Service
}

func main() {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.New(config.Logger)
	logger.SetGlobalLogger(config.Logger)

	appLogger.Infof("Configuration loaded - Port: %d", config.APIPort)

	// Initialize database
	ctx := context.Background()

	database, err := db.NewPostgresDB(ctx, config.Database, appLogger)
	if err != nil {
		appLogger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := migrations.RunMigrations(ctx, database.Pool, appLogger); err != nil {
		appLogger.Fatalf("Failed to run migrations: %v", err)
	}

	// Create application instance
	app := &App{
		config: config,
		logger: appLogger,
		db:     database,
	}

	// Setup auth service
	app.setupAuthService()

	// Setup HTTP server
	if err := app.setupServer(); err != nil {
		appLogger.Fatalf("Failed to setup server: %v", err)
	}

	// Start server with graceful shutdown
	if err := app.start(); err != nil {
		app.logger.Fatalf("Server error: %v", err)
	}
}

// setupServer configures the HTTP server and routes
func (app *App) setupServer() error {
	// Create HTTP mux for routing
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", app.healthCheckHandler)

	// Auth endpoints
	mux.HandleFunc("POST /api/auth/signup", handlers.SignupHandler(app.db, app.logger))

	// Mount auth service routes (auth handler and avatar handler)
	authHandler, avatarHandler := app.auth.Handlers()
	mux.Handle("/auth/", http.StripPrefix("/auth", authHandler))
	mux.Handle("/avatar/", http.StripPrefix("/avatar", avatarHandler))

	// Configure server
	app.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.config.Host, app.config.APIPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

// setupAuthService configures the authentication service
func (app *App) setupAuthService() {
	// Setup auth options
	authOptions := authpkg.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return app.config.Auth.JWTSecret, nil
		}),
		TokenDuration:  time.Duration(app.config.Auth.TokenDuration) * time.Minute,
		CookieDuration: time.Duration(app.config.Auth.CookieDuration) * time.Hour,
		Issuer:         "ferna",
		AudienceReader: token.AudienceFunc(func() ([]string, error) {
			return []string{"ferna-mobile"}, nil
		}),
		URL:         app.config.Auth.BaseURL,
		DisableXSRF: app.config.Auth.DisableXSRF,
		AvatarStore: avatar.NewLocalFS(app.config.Auth.AvatarPath),
	}

	app.auth = authpkg.NewService(authOptions)

	// Add local provider for credential checking
	app.auth.AddDirectProvider("local", provider.CredCheckerFunc(func(user, passwd string) (ok bool, err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get user from database
		dbUser, err := app.db.GetUserByEmail(ctx, user)
		if err != nil {
			app.logger.Debugf("Error getting user for auth: %v", err)
			return false, err
		}
		if dbUser == nil {
			app.logger.Debugf("User not found for auth: %s", user)
			return false, nil
		}

		// Check if user is local auth provider and has password
		if dbUser.AuthProvider != "local" || dbUser.PasswordHash == nil {
			app.logger.Debugf("User is not local auth provider: %s", user)
			return false, nil
		}

		// Verify password using internal auth package
		valid, err := auth.VerifyPassword(passwd, *dbUser.PasswordHash)
		if err != nil {
			app.logger.Debugf("Password verification error: %v", err)
			return false, err
		}

		if valid {
			app.logger.Infof("User authenticated successfully: %s", user)
		} else {
			app.logger.Debugf("Invalid password for user: %s", user)
		}

		return valid, nil
	}))
}

// start starts the server with graceful shutdown
func (app *App) start() error {
	// Channel to listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		app.logger.Infof("Server starting on %s", app.server.Addr)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Errorf("Server failed to start: %v", err)
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	<-quit
	app.logger.Info("Server shutdown signal received...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// handle other cleanup tasks here

	// Attempt graceful shutdown
	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	app.logger.Info("Server exited")
	return nil
}

// HTTP Handlers

// healthCheckHandler provides a simple health check endpoint
func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(`{
		"status": "healthy",
		"service": "%s",
		"timestamp": "%s"
	}`, app.config.Logger.Service, time.Now().UTC().Format(time.RFC3339))

	w.Write([]byte(response))

	app.logger.Debug("Health check accessed")
}
