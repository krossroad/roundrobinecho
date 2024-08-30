package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/krossroad/roundrobinecho/internal/env"
	"github.com/krossroad/roundrobinecho/internal/middlewares"
)

func main() {
	address := env.MustGet("HTTP_ADDRESS")
	logger := bootLogger()
	app := &App{Logger: logger}
	logger.Info("starting echo server", "address", address)

	mux := http.NewServeMux()
	mux.Handle("/echo", middlewares.JSONContentTypeValidator(app.handleEcho))

	if err := http.ListenAndServe(address, mux); err != nil {
		logger.Error("server error", "error", err)
	}
}

func bootLogger() *slog.Logger {
	const version = "1.0.0-beta1"

	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(h)
	return logger.With(
		slog.Group(
			"app-info",
			slog.String("version", version),
			slog.String("server-id", os.Getenv("HOSTNAME")),
		),
	)
}
