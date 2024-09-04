// package echo is the package that contains the echo service.
package echo_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"testing"

	"github.com/krossroad/roundrobinecho/internal/echo"
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
