package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/krossroad/roundrobinecho/internal/echo"
)

type App struct {
	logger  *slog.Logger
	echoSvc echo.Service
}

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
		a.logger.Error("echoSvc.Echo() failed", "error", err)
		a.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

// Error is a helper method that writes an error response.
func (a *App) Error(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	respEnc := json.NewEncoder(w)
	respEnc.Encode(map[string]interface{}{"error": message}) //nolint:errcheck
}
