package main

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krossroad/roundrobinecho/test/mocks"
	"github.com/stretchr/testify/mock"
)

func TestApp_fanoutHandler(t *testing.T) {
	echoSvc := new(mocks.Service)
	tests := []struct {
		name           string
		wantErr        bool
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "case-1/echo-failed",
			mock: func() {
				echoSvc.On("Echo", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "case-2/echo-success",
			mock: func() {
				echoSvc.On("Echo", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	var app *App

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoSvc = new(mocks.Service)
			app = NewApp(slog.Default(), echoSvc)
			req := httptest.NewRequest("GET", "/echo", nil)
			rr := httptest.NewRecorder()

			tt.mock()
			app.fanoutHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rr.Code)
				return
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected response body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
