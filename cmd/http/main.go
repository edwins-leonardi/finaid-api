package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/edwins-leonardi/finaid-api/internal/adapter/config"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/handler/http"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/logger"
)

func main() {
	// Load App configuration
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	logger.Set(config.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
