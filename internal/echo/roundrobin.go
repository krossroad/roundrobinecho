package echo

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

type roundRobin struct {
	log      *slog.Logger
	mu       sync.Mutex
	backends []*Backend
	counter  atomic.Uint32
}

var _ LoadBalancer = (*roundRobin)(nil)

func NewRoundRobin(log *slog.Logger, startingAddresses []string) (LoadBalancer, error) {
	backends := []*Backend{}

	for _, v := range startingAddresses {
		b, err := NewBackend(v)
		if err != nil {
			return nil, fmt.Errorf("NewBackend(): %w", err)
		}

		backends = append(backends, b)
	}

	return &roundRobin{
		backends: backends,
		log:      log,
	}, nil
}

func (b *roundRobin) Backends() []*Backend {
	return b.backends
}

func (rr *roundRobin) Next() (*Backend, error) {
	backendLen := len(rr.backends)
	attempts := 0

	rr.mu.Lock()
	defer rr.mu.Unlock()
	for {
		if attempts > backendLen {
			return nil, fmt.Errorf("all backends are down")
		}

		attempts++

		i := rr.counter.Add(1) % uint32(backendLen)
		if rr.backends[i].Alive() {
			return rr.backends[i], nil
		}
	}
}
