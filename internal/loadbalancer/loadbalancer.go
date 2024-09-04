package loadbalancer

import (
	"net/http"
	"net/url"
)

// Backend represents a backend server in the load balancer.
type Backend interface {
	// Alive returns true if the backend server is alive, false otherwise.
	Alive() bool

	// SetAlive sets the alive status of the backend server.
	SetAlive(bool)

	// Do handles the incoming HTTP request using the backend server.
	Do(http.ResponseWriter, *http.Request)

	// HealthCheckURL returns the URL used for health checks of the backend server.
	HealthCheckURL() string

	// URL returns the clone of URL of the backend server.
	URL() *url.URL
}

// LoadBalancer represents a load balancer interface.
type LoadBalancer interface {
	// Next returns the next available backend for load balancing.
	// It returns an error if there are no available backends.
	Next() (Backend, error)

	// Backends returns a slice of all the available backends.
	Backends() []Backend
}
