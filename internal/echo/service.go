// package echo is the package that contains the echo service.
package echo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/krossroad/roundrobinecho/internal/loadbalancer"
)

type (
	// Service is the echo service.
	Service interface {
		Echo(ctx context.Context, w http.ResponseWriter, r *http.Request) error
		Monitor(ctx context.Context)
	}

	echoService struct {
		logger *slog.Logger
		loadbalancer.LoadBalancer
		healthCheckInterval time.Duration
	}

	// OptSetter is a function that sets an option on the Service.
	OptSetter func(*echoService)
)

var _ Service = (*echoService)(nil)

// NewService creates a new instance of the echo Service.
// It takes a logger, a load balancer, and an optional list of OptSetters as parameters.
func NewService(log *slog.Logger, lb loadbalancer.LoadBalancer, setter ...OptSetter) *echoService {
	svc := &echoService{
		logger:       log,
		LoadBalancer: lb,
	}

	for _, optSetter := range setter {
		optSetter(svc)
	}

	return svc
}

// WithHealthCheckInterval sets the health check interval for the echo service.
func WithHealthCheckInterval(interval time.Duration) OptSetter {
	return func(svc *echoService) {
		svc.healthCheckInterval = interval
	}
}

// Echo forwards the request to the next service in the round-robin list.
func (svc *echoService) Echo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	backend, err := svc.Next()
	if err != nil {
		return fmt.Errorf("Next(): %w", err)
	}
	svc.logger.Info("forwarding request", "backend", backend.URL().String())
	backend.Do(w, r)

	return nil
}

// Monitor is responsible for periodically checking the health
// status of each backend server in the load balancer.
func (svc *echoService) Monitor(ctx context.Context) {
	timer := time.NewTicker(svc.healthCheckInterval)
	for {
		for _, backend := range svc.Backends() {
			_, err := http.Get(backend.HealthCheckURL())
			if err != nil {
				backend.SetAlive(false)
			} else {
				backend.SetAlive(true)
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
	}
}
