package openai

import (
	"errors"
	"strings"
)

var ErrTokenRequired = errors.New("openai token required")

type Config struct {
	Token string
}

func NewConfig(token string) *Config {
	return &Config{Token: token}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return ErrTokenRequired
	}

	return nil
}
