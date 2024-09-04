package echo

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/krossroad/roundrobinecho/internal/loadbalancer"
)

type echoBackend struct {
	url          *url.URL
	alive        bool
	mu           sync.RWMutex
	reverseProxy *httputil.ReverseProxy
}

var _ loadbalancer.Backend = (*echoBackend)(nil)

// NewBackend creates a new instance of the echoBackend struct.
// It takes a target address as a parameter and returns a pointer to the echoBackend struct.
// The echoBackend struct implements the Backend interface and is responsible for handling requests to the backend service.
func NewBackend(target string) (loadbalancer.Backend, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %w", err)
	}

	return &echoBackend{
		url:          u,
		reverseProxy: httputil.NewSingleHostReverseProxy(u),
	}, nil
}

// Alive returns a boolean value indicating whether the backend is alive or not.
func (b *echoBackend) Alive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

// SetAlive sets the alive status of the backend.
func (b *echoBackend) SetAlive(alive bool) {
	b.mu.Lock()
	b.alive = alive
	b.mu.Unlock()
}

// Do forwards the request to the backend service.
func (b *echoBackend) Do(w http.ResponseWriter, r *http.Request) {
	b.reverseProxy.ServeHTTP(w, r)
}

// HealthCheckURL returns the health check URL for the backend service.
func (b *echoBackend) HealthCheckURL() string {
	u := b.URL()
	u.Path = "/health"

	return u.String()
}

// URL returns the URL of the backend service.
func (b *echoBackend) URL() (u *url.URL) {
	u, _ = url.Parse(b.url.String())
	return
}
