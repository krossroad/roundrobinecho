package main

import (
	"log/slog"
	"net/http"

	"github.com/krossroad/roundrobinecho/internal/echo"
)

type (
	App struct {
		logger  *slog.Logger
		echoSvc echo.Service
	}
)

// NewApp creates a new App.
// It takes a logger and an echo service as arguments.
// It returns a new App.
func NewApp(logger *slog.Logger, echoSvc echo.Service) *App {
	return &App{
		logger:  logger,
		echoSvc: echoSvc,
	}
}

// fanoutHandler is the HTTP handler for the round-robin endpoint.
func (a *App) fanoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.echoSvc.Echo(r.Context(), w, r); err != nil {
		a.logger.Error("echo error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
