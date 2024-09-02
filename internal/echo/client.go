// package echo is the package that contains the echo service.
package echo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type (
	Backend struct {
		url          *url.URL
		alive        bool
		mu           sync.RWMutex
		reverseProxy *httputil.ReverseProxy
	}
	LoadBalancer interface {
		Next() (*Backend, error)
		Backends() []*Backend
	}

	// Service is the echo service.
	Service struct {
		LoadBalancer
		logger *slog.Logger
	}
)

func NewBackend(target string) (*Backend, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %w", err)
	}

	return &Backend{
		url:          url,
		reverseProxy: httputil.NewSingleHostReverseProxy(url),
	}, nil
}

// Alive returns a boolean value indicating whether the backend is alive or not.
func (b *Backend) Alive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.alive = alive
	b.mu.Unlock()
}

func (b *Backend) Forward(w http.ResponseWriter, r *http.Request) {
	b.reverseProxy.ServeHTTP(w, r)
}

func (b *Backend) HealthCheckURL() string {
	u, _ := url.Parse(b.url.String())
	u.Path = "/health"

	return u.String()
}

// NewService creates a new echo service..
// It takes a logger and a list of service addresses as input parameters.
// It returns a pointer to the Service struct and an error.
func NewService(log *slog.Logger, lb LoadBalancer) (*Service, error) {
	svc := &Service{
		logger:       log,
		LoadBalancer: lb,
	}

	// Health-check
	go svc.HealthCheck()

	return svc, nil
}

// Echo forwards the request to the next service in the round-robin list.
func (svc *Service) Echo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	backend, err := svc.Next()
	if err != nil {
		return fmt.Errorf("Next(): %w", err)
	}
	svc.logger.Info("forwarding request", "backend", backend.url.String())
	backend.Forward(w, r)

	return nil
}

func (svc *Service) HealthCheck() {
	timer := time.NewTicker(30 * time.Second)
	for {
		for _, backend := range svc.Backends() {
			_, err := http.Get(backend.HealthCheckURL())
			if err != nil {
				backend.SetAlive(false)
			} else {
				backend.SetAlive(true)
			}
		}

		<-timer.C
	}
}
