package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRoundTripper is a mock implementation of http.RoundTripper.
type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

// RoundTrip intercepts the request and returns a mock response or error.
func (m *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Token: "test-token",
			},
			expectError: false,
		},
		{
			name: "empty token",
			config: &Config{
				Token: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.config.Token, service.token)
				assert.Equal(t, "https://api.replicate.com/v1/models/black-forest-labs/flux-1.1-pro-ultra/predictions", service.url)
			}
		})
	}
}

func TestService_GenerateImage(t *testing.T) {
	tests := []struct {
		name           string
		request        *Request
		mockResponse   *Response
		mockStatusCode int
		mockErr        error
		expectError    bool
		expectedError  error
	}{
		{
			name:    "successful generation",
			request: &Request{
				// заполните необходимыми данными
			},
			mockResponse: &Response{
				ID:     "test-id",
				Status: "completed",
				Output: "test-output",
			},
			mockStatusCode: http.StatusCreated,
			expectError:    false,
		},
		{
			name:    "bad request",
			request: &Request{
				// заполните необходимыми данными
			},
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
			expectedError:  ErrBadRequest,
		},
		{
			name:        "network error",
			request:     &Request{},
			mockErr:     io.EOF,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготавливаем моковый ответ
			var respBody []byte
			var err error
			if tt.mockResponse != nil {
				respBody, err = json.Marshal(tt.mockResponse)
				require.NoError(t, err)
			}

			// Создаем моковый транспорт
			mockTransport := &MockRoundTripper{
				Response: &http.Response{
					StatusCode: tt.mockStatusCode,
					Body:       io.NopCloser(bytes.NewBuffer(respBody)),
				},
				Err: tt.mockErr,
			}

			// Создаем HTTP клиент с моковым транспортом
			originalClient := http.DefaultClient
			http.DefaultClient = &http.Client{
				Transport: mockTransport,
			}
			defer func() {
				http.DefaultClient = originalClient
			}()

			service := &Service{
				token: "test-token",
				url:   "test-url",
			}

			resp, err := service.GenerateImage(context.Background(), tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.ID, resp.ID)
				assert.Equal(t, tt.mockResponse.Status, resp.Status)
				assert.Equal(t, tt.mockResponse.Output, resp.Output)
			}
		})
	}
}

func TestService_GenerateImage_InvalidJSON(t *testing.T) {
	mockTransport := &MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
		},
	}

	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: mockTransport,
	}
	defer func() {
		http.DefaultClient = originalClient
	}()

	service := &Service{
		token: "test-token",
		url:   "test-url",
	}

	_, err := service.GenerateImage(context.Background(), &Request{})
	assert.Error(t, err)
}
