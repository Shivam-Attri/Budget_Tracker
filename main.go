// main.go
// The application entry point. It initializes all components and starts the server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/config"
	"your_username/budget-tracker/database"
	"your_username/budget-tracker/handlers"
	"your_username/budget-tracker/middleware"
	"your_username/budget-tracker/router"
	"your_username/budget-tracker/utils"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	// Setup structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration using Viper
	cfg, err := config.LoadConfig(".")
	if err != nil {
		zlog.Fatal().Err(err).Msg("Could not load configuration")
	}

	// Connect to the database
	db := database.ConnectDB(cfg.DB.URL)
	defer db.Close()

	// Create/update the necessary tables
	database.CreateTables(db)

	// Initialize packages
	auth.Init(cfg.Server.JWTSecret, cfg.Server.AccessTokenTTL, cfg.Server.RefreshTokenTTL)
	utils.InitValidator()
	middleware.InitRateLimiter(cfg.RateLimiter)

	// Create the data store
	store := &database.DBStore{
		DB:  db,
		Key: []byte(cfg.Server.EncryptionKey),
	}

	// Create the environment for dependency injection
	env := &handlers.Env{
		DB:                store,
		DefaultPageLimit:  cfg.Server.DefaultPageLimit,
		AllowedSortFields: cfg.Server.AllowedSortFields,
	}

	// Initialize the router
	r := router.NewRouter(env)

	// Setup the HTTP server
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		WriteTimeout: time.Duration(cfg.Server.Timeout.Write) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.Timeout.Read) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Timeout.Idle) * time.Second,
	}

	// Run server in a goroutine
	go func() {
		zlog.Info().Str("address", srv.Addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.Timeout.GracefulShutdown)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	zlog.Info().Msg("Server exiting")
}
