package echo

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchoBackend_Alive(t *testing.T) {
	tests := []struct {
		name  string
		alive bool
	}{
		{
			name:  "case-1/Backend is alive",
			alive: true,
		},
		{
			name:  "case-2/Backend is not alive",
			alive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &echoBackend{
				alive: tt.alive,
			}

			got := b.Alive()
			if got != tt.alive {
				t.Errorf("Alive() = %v, want %v", got, tt.alive)
			}
		})
	}
}

func TestEchoBackend_SetAlive(t *testing.T) {
	tests := []struct {
		name  string
		alive bool
	}{
		{
			name:  "case-1/Set backend alive",
			alive: true,
		},
		{
			name:  "case-2/Set backend not alive",
			alive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &echoBackend{
				alive: !tt.alive,
			}

			b.SetAlive(tt.alive)

			got := b.Alive()
			if got != tt.alive {
				t.Errorf("SetAlive() = %v, want %v", got, tt.alive)
			}
		})
	}
}



func TestEchoBackend_HealthCheckURL(t *testing.T) {
	uStr := "http://example.com/health"
	u, _ := url.Parse("http://example.com")

	t.Run("case-1/valid-URL", func(t *testing.T) {
		b := &echoBackend{
			url: u,
		}

		got := b.HealthCheckURL()
		if got != uStr {
			t.Errorf("HealthCheckURL() = %v, want %v", got, uStr)
		}
	})

}

func TestEchoBackend_URL(t *testing.T) {
	u, _ := url.Parse("http://example.com")
	t.Run("case-1/valid", func(t *testing.T) {
		b := &echoBackend{
			url: u,
		}

		got := b.URL()
		assert.Equal(t, u.String(), got.String())
	})
}
