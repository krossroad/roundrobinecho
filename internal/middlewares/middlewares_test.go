package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krossroad/roundrobinecho/internal/middlewares"
)

func TestJSONContentTypeValidator(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		expectedCode int
	}{
		{
			name:         "case-1/Valid content type",
			contentType:  "application/json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case-2/Invalid content type",
			contentType:  "text/plain",
			expectedCode: http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middlewares.JSONContentTypeValidator(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v, want %v", rr.Code, tt.expectedCode)
			}
		})
	}
}
