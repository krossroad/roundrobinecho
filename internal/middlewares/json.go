package middlewares

import (
	"encoding/json"
	"net/http"
)

func JSONContentTypeValidator(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			writeError(w, http.StatusUnsupportedMediaType, "unsupported content type")
			return
		}

		fn(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	respEnc := json.NewEncoder(w)
	respEnc.Encode(map[string]string{"error": message})
}
