package env_test

import (
	"os"
	"testing"

	"github.com/krossroad/roundrobinecho/internal/env"
)

func TestMustGet(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		expectedPanic string
	}{
		{
			name:          "case-1/Existing environment variable",
			key:           "EXISTING_VAR",
			value:         "existing_value",
			expectedPanic: "",
		},
		{
			name:          "case-2/Missing environment variable",
			key:           "MISSING_VAR",
			value:         "",
			expectedPanic: "missing required environment variable: MISSING_VAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.key, tt.value)

			defer func() {
				r := recover()
				if r != nil {
					if r != tt.expectedPanic {
						t.Errorf("MustGet() panic = %v, expected %v", r, tt.expectedPanic)
					}
				} else {
					if tt.expectedPanic != "" {
						t.Errorf("MustGet() did not panic, expected panic: %v", tt.expectedPanic)
					}
				}
			}()

			env.MustGet(tt.key)
		})
	}
}
