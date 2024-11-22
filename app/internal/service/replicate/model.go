package replicate

import (
	"math/rand"
	"time"
)

const (
	DefaultRatio               = AspectRatio16x9
	DefaultRaw                 = false
	DefaultSafetyTolerance     = 4
	DefaultImagePromptStrength = 0.3
)

type Request struct {
	Input *Input `json:"input"`
}

func NewRequest(options ...Option) *Request {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	defaultSeed := rng.Uint64()

	// Create a default Input instance with default values.
	input := &Input{
		Ratio:               DefaultRatio.String(),
		Raw:                 DefaultRaw,
		SafetyTolerance:     DefaultSafetyTolerance,
		Seed:                defaultSeed,
		ImagePromptStrength: DefaultImagePromptStrength,
	}

	// Apply all provided options to configure the instance.
	for _, opt := range options {
		opt(input)
	}

	return &Request{
		Input: input,
	}
}

// Input represents the input structure with configurable fields.
type Input struct {
	Prompt              string  `json:"prompt"`
	Ratio               string  `json:"aspect_ratio"`
	Raw                 bool    `json:"raw"`
	SafetyTolerance     uint8   `json:"safety_tolerance"`
	Seed                uint64  `json:"seed"`
	ImagePromptStrength float64 `json:"image_prompt_strength"`
}

// Option defines a functional option for configuring an Input instance.
type Option func(*Input)

// WithPrompt sets the Prompt field.
func WithPrompt(prompt string) Option {
	return func(i *Input) {
		i.Prompt = prompt
	}
}

// WithRatio sets the Ratio field.
func WithRatio(ratio AspectRatio) Option {
	return func(i *Input) {
		i.Ratio = ratio.String()
	}
}

// WithRaw sets the Raw field.
func WithRaw(raw bool) Option {
	return func(i *Input) {
		i.Raw = raw
	}
}

// WithSafetyTolerance sets the SafetyTolerance field.
func WithSafetyTolerance(tolerance uint8) Option {
	return func(i *Input) {
		i.SafetyTolerance = tolerance
	}
}

// WithSeed sets the Seed field explicitly.
func WithSeed(seed uint64) Option {
	return func(i *Input) {
		i.Seed = seed
	}
}

// WithImagePromptStrength sets the ImagePromptStrength field.
func WithImagePromptStrength(strength float64) Option {
	return func(i *Input) {
		i.ImagePromptStrength = strength
	}
}

type Response struct {
	ID          string      `json:"id"`
	Model       string      `json:"model"`
	Version     string      `json:"version"`
	Input       Input       `json:"input"`
	Logs        string      `json:"logs"`
	Output      string      `json:"output"`
	DataRemoved bool        `json:"data_removed"`
	Error       interface{} `json:"error"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
}
