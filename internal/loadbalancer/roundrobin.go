package loadbalancer

import (
	"fmt"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
)

type roundRobin struct {
	log      *slog.Logger
	mu       sync.Mutex
	backends []Backend
	counter  atomic.Uint32
}

var _ LoadBalancer = (*roundRobin)(nil)

// NewRoundRobin creates a new round-robin load balancer with the given logger and starting addresses.
// It returns a LoadBalancer interface and an error if there are no backends provided or if there is an error creating a backend.
// The LoadBalancer interface allows for distributing requests among multiple backends in a round-robin fashion.
func NewRoundRobin(log *slog.Logger, backends []Backend) (LoadBalancer, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends provided")
	}

	return &roundRobin{
		backends: backends,
		log:      log,
	}, nil
}

func (rr *roundRobin) Backends() []Backend {
	return rr.backends
}

// Next returns the next available backend in a round-robin manner.
// It selects a backend from the list of backends based on the current counter value.
// If all backends are down, it returns an error.
func (rr *roundRobin) Next() (Backend, error) {
	if len(rr.backends) > math.MaxUint32 {
		return nil, fmt.Errorf("too many backends")
	}

	attempts := uint32(0)
	backendLen := uint32(len(rr.backends)) //nolint:gosec
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for {
		if attempts > backendLen {
			return nil, fmt.Errorf("all backends are down")
		}

		attempts++

		i := rr.counter.Add(1) % backendLen
		if rr.backends[i].Alive() {
			return rr.backends[i], nil
		}
	}
}
