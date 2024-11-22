package replicate

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewConfig(t *testing.T) {
	t.Run("WithToken", func(t *testing.T) {
		token := "test-token"
		config := NewConfig(token)

		if config.Token != token {
			t.Errorf("expected Token to be %q, got %q", token, config.Token)
		}
	})

	t.Run("WithoutToken", func(t *testing.T) {
		config := NewConfig("")

		if config.Token != "" {
			t.Errorf("expected Token to be empty, got %q", config.Token)
		}
	})
}

func TestConfig_Validate(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		config := NewConfig("valid-token")

		err := config.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("MissingToken", func(t *testing.T) {
		config := NewConfig("")

		err := config.Validate()
		if err == nil {
			t.Errorf("expected error %v, got nil", ErrTokenRequired)
		}

		if !errors.Is(err, ErrTokenRequired) {
			t.Errorf("expected error %v, got %v", ErrTokenRequired, err)
		}
	})
}
