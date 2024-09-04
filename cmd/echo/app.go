package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type App struct {
	*slog.Logger
}

func (a *App) handleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if r.Body == nil {
		a.Error(w, http.StatusBadRequest, "empty request body")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Error(w, http.StatusInternalServerError, "error reading request body")
		return
	}
	defer r.Body.Close()

	var resp json.RawMessage
	if err := json.Unmarshal(body, &resp); err != nil {
		a.Error(w, http.StatusBadRequest, "error parsing request body")
		return
	}

	if _, err := w.Write(body); err != nil {
		a.Error(w, http.StatusInternalServerError, "error writing response")
		return
	}

	a.Debug("request echoed")
}

func (a *App) Error(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	respEnc := json.NewEncoder(w)
	respEnc.Encode(map[string]interface{}{"error": message})
}
