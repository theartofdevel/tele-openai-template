package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected *Config
	}{
		{
			name:     "with valid token",
			token:    "sk-1234567890",
			expected: &Config{Token: "sk-1234567890"},
		},
		{
			name:     "with empty token",
			token:    "",
			expected: &Config{Token: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.token)
			assert.Equal(t, tt.expected, config)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError error
	}{
		{
			name:      "valid config",
			config:    &Config{Token: "sk-1234567890"},
			wantError: nil,
		},
		{
			name:      "empty token",
			config:    &Config{Token: ""},
			wantError: ErrTokenRequired,
		},
		{
			name:      "whitespace token",
			config:    &Config{Token: "   "},
			wantError: ErrTokenRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
