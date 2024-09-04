package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"

	"github.com/krossroad/roundrobinecho/internal/env"
	"github.com/krossroad/roundrobinecho/internal/logger"
	"github.com/krossroad/roundrobinecho/internal/middlewares"
)

func main() {
	address := env.MustGet("HTTP_ADDRESS")
	logger := logger.New("echo")
	app := &App{Logger: logger}
	logger.Info("starting echo server", "address", address)

	mux := http.NewServeMux()
	mux.Handle("/echo", middlewares.JSONContentTypeValidator(app.handleEcho))
	mux.HandleFunc("/health", bootHealthCheck())

	srv := &http.Server{
		Addr:        address,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "error", err)
	}
}

func bootHealthCheck() http.HandlerFunc {
	checker := health.NewChecker(
		health.WithPeriodicCheck(5*time.Second, 0*time.Second, health.Check{
			Name: "echo",
			Check: func(_ context.Context) error {
				return nil
			},
		}),
	)

	return health.NewHandler(checker)
}
