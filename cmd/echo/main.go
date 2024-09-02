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

	if err := http.ListenAndServe(address, mux); err != nil {
		logger.Error("server error", "error", err)
	}
}

func bootHealthCheck() http.HandlerFunc {
	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithPeriodicCheck(7*time.Second, 5*time.Second, health.Check{
			Name: "echo",
			Check: func(ctx context.Context) error {
				return nil
			},
		}),
	)

	return health.NewHandler(checker)
}
