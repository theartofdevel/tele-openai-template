package bot

import (
	"errors"
	"strings"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

var (
	ErrTokenRequired = errors.New("bot token required")
)

type Config struct {
	Token     string
	Timeout   time.Duration
	Whitelist []int64
}

func NewConfig(token string, timeout time.Duration, whitelist []int64) *Config {
	return &Config{
		Token:     token,
		Timeout:   timeout,
		Whitelist: whitelist,
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return ErrTokenRequired
	}

	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}

	return nil
}
