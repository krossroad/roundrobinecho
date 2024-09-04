package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApp_handleEcho(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           io.Reader
		expectedStatus int
		expectedBody   string
		ContentType    string
	}{
		{
			name:           "BadRequest/EmptyBody",
			method:         http.MethodPost,
			expectedBody:   `{"error":"empty request body"}` + "\n",
			expectedStatus: http.StatusBadRequest,
			ContentType:    "application/json",
		},
		{
			name:           "BadRequest/InvalidBody",
			expectedBody:   `{"error":"error parsing request body"}` + "\n",
			method:         http.MethodPost,
			body:           strings.NewReader(`{"foo":"bar`),
			expectedStatus: http.StatusBadRequest,
			ContentType:    "application/json",
		},
		{
			name:           "Ok",
			method:         http.MethodPost,
			body:           strings.NewReader(`{"foo":"bar"}`),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"foo":"bar"}`,
			ContentType:    "application/json",
		},
		{
			name:           "MethodNotAllowed",
			method:         http.MethodGet,
			expectedBody:   `{"error":"method not allowed"}` + "\n",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	app := &App{
		Logger: slog.Default(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/echo", tt.body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.ContentType)

			rr := httptest.NewRecorder()

			app.handleEcho(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rr.Code)
				return
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected response body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
