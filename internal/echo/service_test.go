// package echo is the package that contains the echo service.
package echo_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/krossroad/roundrobinecho/internal/echo"
	"github.com/krossroad/roundrobinecho/internal/loadbalancer"
	"github.com/krossroad/roundrobinecho/test/mocks"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
)

func TestService_Echo(t *testing.T) {
	rr := new(mocks.LoadBalancer)
	tests := []struct {
		name    string
		wantErr error
		mocked  func()
	}{
		{
			name: "case-1/next-failed",
			mocked: func() {
				rr.On("Next").Return(nil, fmt.Errorf("random error"))
			},
			wantErr: errors.New("Next(): random error"),
		},

		{
			name: "case-2/next-success",
			mocked: func() {
				b := new(mocks.Backend)
				b.On("URL").Return(&url.URL{})
				b.On("Do", m.Anything, m.Anything)

				rr.On("Next").Return(b, nil)
			},
		},
	}

	log := slog.Default()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr = new(mocks.LoadBalancer)
			if tt.mocked != nil {
				tt.mocked()
			}
			svc := echo.NewService(log, rr)
			err := svc.Echo(context.Background(), nil, nil)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestService_Monitor(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.Default()

	// Mock the health check for the backends
	b1 := new(mocks.Backend)
	b1.On("HealthCheckURL").Return("http://backend1/health").Once()
	b1.On("SetAlive", true).Once()

	b2 := new(mocks.Backend)
	b2.On("HealthCheckURL").Return("http://backend2/health").Once()
	b2.On("SetAlive", false).Once()

	lb := new(mocks.LoadBalancer)
	lb.On("Backends").Return([]loadbalancer.Backend{b1, b2})

	resp1 := &http.Response{
		StatusCode: http.StatusOK,
	}
	resp2 := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
	}
	trMock := new(mocks.RoundTripper)
	trMock.On("RoundTrip", m.Anything).Return(resp1, nil).Once()
	trMock.On("RoundTrip", m.Anything).Return(resp2, nil).Once()
	client := &http.Client{
		Transport: trMock,
	}

	svc := echo.NewService(log, lb, echo.WithHealthCheckInterval(10*time.Second), echo.WithHTTPClient(client))

	// Mock the responses for the health check requests
	// Start monitoring
	go svc.Monitor(ctx)
	time.Sleep(2 * time.Second)
}
