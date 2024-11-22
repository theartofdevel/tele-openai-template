package replicate

import (
	"errors"
)

var ErrTokenRequired = errors.New("replicate token required")

type Config struct {
	Token string
}

func NewConfig(token string) *Config {
	return &Config{Token: token}
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return ErrTokenRequired
	}

	return nil
}
