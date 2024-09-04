package loadbalancer_test

import (
	"fmt"
	"log/slog"
	"reflect"
	"testing"

	"github.com/krossroad/roundrobinecho/internal/echo"
	lb "github.com/krossroad/roundrobinecho/internal/loadbalancer"
)

func TestRoundRobinNextAllBackendDead(t *testing.T) {
	log := slog.Default()
	b1, _ := echo.NewBackend("backend1")
	backends := []lb.Backend{b1}

	rr, err := lb.NewRoundRobin(log, backends)
	if err != nil {
		t.Fatalf("Failed to create roundRobin: %v", err)
	}

	_, err = rr.Next()
	if err == nil {
		t.Errorf("RoundRobin.Next() error = %v, wantErr %v", err, true)
	}
}
func TestRoundRobinNext(t *testing.T) {
	log := slog.Default()
	b1, _ := echo.NewBackend("backend1")
	b2, _ := echo.NewBackend("backend2")
	b3, _ := echo.NewBackend("backend3")

	rr, err := lb.NewRoundRobin(log, []lb.Backend{b1, b2, b3})
	if err != nil {
		t.Fatalf("Failed to create roundRobin: %v", err)
	}

	rr.Backends()[0].SetAlive(true)
	rr.Backends()[2].SetAlive(true)

	// Call Next multiple times and verify the returned backend
	tests := []struct {
		name     string
		expected lb.Backend
		wantErr  bool
	}{
		{
			name:     "First call",
			expected: rr.Backends()[2],
			wantErr:  false,
		},
		{
			name:     "Second call",
			expected: rr.Backends()[0],
			wantErr:  false,
		},
		{
			name:     "Third call",
			expected: rr.Backends()[2],
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rr.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundRobin.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RoundRobin.Next() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestNewRoundRobin(t *testing.T) {
	log := slog.Default()
	b1, _ := echo.NewBackend("backend1")
	b2, _ := echo.NewBackend("backend2")
	b3, _ := echo.NewBackend("backend3")

	tests := []struct {
		name          string
		args          []lb.Backend
		expectedError error
	}{
		{
			name:          "Empty starting addresses",
			args:          []lb.Backend{},
			expectedError: fmt.Errorf("no backends provided"),
		},
		{
			name:          "Valid starting addresses",
			args:          []lb.Backend{b1, b2, b3},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := lb.NewRoundRobin(log, tt.args)
			if tt.expectedError != nil && (err == nil || err.Error() != tt.expectedError.Error()) {
				t.Errorf("NewRoundRobin() error = %q, expected %q", err, tt.expectedError)
			}
		})
	}
}
