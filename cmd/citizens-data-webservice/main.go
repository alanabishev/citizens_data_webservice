// Package main provides the entry point to the citizen_webservice application.
package main

import (
	"citizen_webservice/internal/config"
	handlerDelete "citizen_webservice/internal/http-server/handlers/delete"
	"citizen_webservice/internal/http-server/handlers/get"
	"citizen_webservice/internal/http-server/handlers/health"
	"citizen_webservice/internal/http-server/handlers/iin_validate"
	"citizen_webservice/internal/http-server/handlers/save"
	"citizen_webservice/internal/storage/sqlite"

	mwLogger "citizen_webservice/internal/http-server/middleware/logger"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// main is the entry point of the application.
// It sets up the configuration, logger, storage, and HTTP server.
// It also handles graceful shutdown of the server on receiving an interrupt signal.
func main() {
	// 1. Config
	cfg := config.Load()
	fmt.Println(cfg)

	// 2. Logger
	log := setupLogger(cfg.Env)
	log.Info(
		"starting citizens-data-webservice",
		slog.String("env", cfg.Env),
		slog.String("version", "1.0.0"),
	)
	log.Debug("debug messages are enabled")

	// 3. Storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", slog.Any("error", err))
		os.Exit(1)
	}

	// 4. Router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.New(log))

	// Health check endpoint (no auth required)
	router.Get("/health", health.Check(log))

	// Define the routes for the HTTP server.
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.BasicAuth("citizen_website", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Get("/iin_check/{iin}", iin_validate.Execute(log))
		r.Post("/people/info", save.Person(log, storage))
		r.Get("/people/info/iin/{iin}", get.ByIIN(log, storage))
		r.Get("/people/info/name/{name}", get.ByName(log, storage))
		r.Delete("/people/delete/{iin}", handlerDelete.ByIIN(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Start the HTTP server in a separate goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", slog.Any("error", err))
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully.
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server")
		return
	}

	log.Info("server stopped")

}

// setupLogger initializes a logger based on the environment.
// It returns a logger with different formats and levels for different environments.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
