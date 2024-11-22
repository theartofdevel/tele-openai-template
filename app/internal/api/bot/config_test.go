package bot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		timeout   time.Duration
		whitelist []int64
		expected  *Config
	}{
		{
			name:      "full config",
			token:     "bot123token",
			timeout:   5 * time.Second,
			whitelist: []int64{123, 456},
			expected: &Config{
				Token:     "bot123token",
				Timeout:   5 * time.Second,
				Whitelist: []int64{123, 456},
			},
		},
		{
			name:      "empty whitelist",
			token:     "bot123token",
			timeout:   3 * time.Second,
			whitelist: []int64{},
			expected: &Config{
				Token:     "bot123token",
				Timeout:   3 * time.Second,
				Whitelist: []int64{},
			},
		},
		{
			name:      "nil whitelist",
			token:     "bot123token",
			timeout:   3 * time.Second,
			whitelist: nil,
			expected: &Config{
				Token:     "bot123token",
				Timeout:   3 * time.Second,
				Whitelist: nil,
			},
		},
		{
			name:      "empty token",
			token:     "",
			timeout:   1 * time.Second,
			whitelist: []int64{123},
			expected: &Config{
				Token:     "",
				Timeout:   1 * time.Second,
				Whitelist: []int64{123},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.token, tt.timeout, tt.whitelist)
			assert.Equal(t, tt.expected, config)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name            string
		config          *Config
		expectedError   error
		expectedTimeout time.Duration
	}{
		{
			name: "valid config",
			config: &Config{
				Token:     "bot123token",
				Timeout:   5 * time.Second,
				Whitelist: []int64{123, 456},
			},
			expectedError:   nil,
			expectedTimeout: 5 * time.Second,
		},
		{
			name: "empty token",
			config: &Config{
				Token:     "",
				Timeout:   5 * time.Second,
				Whitelist: []int64{123},
			},
			expectedError:   ErrTokenRequired,
			expectedTimeout: 5 * time.Second,
		},
		{
			name: "zero timeout",
			config: &Config{
				Token:     "bot123token",
				Timeout:   0,
				Whitelist: []int64{123},
			},
			expectedError:   nil,
			expectedTimeout: defaultTimeout,
		},
		{
			name: "negative timeout",
			config: &Config{
				Token:     "bot123token",
				Timeout:   -1 * time.Second,
				Whitelist: []int64{123},
			},
			expectedError:   nil,
			expectedTimeout: defaultTimeout,
		},
		{
			name: "whitespace token",
			config: &Config{
				Token:     "   ",
				Timeout:   5 * time.Second,
				Whitelist: []int64{123},
			},
			expectedError:   ErrTokenRequired,
			expectedTimeout: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedTimeout, tt.config.Timeout)
		})
	}
}

func TestDefaultTimeout(t *testing.T) {
	assert.Equal(t, 10*time.Second, defaultTimeout)
}
