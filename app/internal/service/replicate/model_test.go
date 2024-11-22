package replicate

import (
	"testing"
)

func TestNewRequest_Defaults(t *testing.T) {
	req := NewRequest()

	// Validate default values
	if req.Input.Ratio != DefaultRatio.String() {
		t.Errorf("expected default Ratio to be %q, got %q", DefaultRatio, req.Input.Ratio)
	}

	if req.Input.Raw != DefaultRaw {
		t.Errorf("expected default Raw to be %v, got %v", DefaultRaw, req.Input.Raw)
	}

	if req.Input.SafetyTolerance != DefaultSafetyTolerance {
		t.Errorf("expected default SafetyTolerance to be %d, got %d", DefaultSafetyTolerance, req.Input.SafetyTolerance)
	}

	if req.Input.ImagePromptStrength != DefaultImagePromptStrength {
		t.Errorf("expected default ImagePromptStrength to be %f, got %f", DefaultImagePromptStrength, req.Input.ImagePromptStrength)
	}

	if req.Input.Seed == 0 {
		t.Errorf("expected default Seed to be non-zero, got %d", req.Input.Seed)
	}

	if req.Input.Prompt != "" {
		t.Errorf("expected default Prompt to be empty, got %q", req.Input.Prompt)
	}
}

func TestNewRequest_WithOptions(t *testing.T) {
	customPrompt := "Custom prompt"
	customRatio := AspectRatio2x3
	customRaw := true
	customTolerance := uint8(2)
	customSeed := uint64(12345)
	customStrength := 0.7

	req := NewRequest(
		WithPrompt(customPrompt),
		WithRatio(customRatio),
		WithRaw(customRaw),
		WithSafetyTolerance(customTolerance),
		WithSeed(customSeed),
		WithImagePromptStrength(customStrength),
	)

	// Validate custom values
	if req.Input.Prompt != customPrompt {
		t.Errorf("expected Prompt to be %q, got %q", customPrompt, req.Input.Prompt)
	}

	if req.Input.Ratio != customRatio.String() {
		t.Errorf("expected Ratio to be %q, got %q", customRatio, req.Input.Ratio)
	}

	if req.Input.Raw != customRaw {
		t.Errorf("expected Raw to be %v, got %v", customRaw, req.Input.Raw)
	}

	if req.Input.SafetyTolerance != customTolerance {
		t.Errorf("expected SafetyTolerance to be %d, got %d", customTolerance, req.Input.SafetyTolerance)
	}

	if req.Input.Seed != customSeed {
		t.Errorf("expected Seed to be %d, got %d", customSeed, req.Input.Seed)
	}

	if req.Input.ImagePromptStrength != customStrength {
		t.Errorf("expected ImagePromptStrength to be %f, got %f", customStrength, req.Input.ImagePromptStrength)
	}
}

func TestNewRequest_WithPartialOptions(t *testing.T) {
	customPrompt := "Partial prompt"
	customTolerance := uint8(10)

	req := NewRequest(
		WithPrompt(customPrompt),
		WithSafetyTolerance(customTolerance),
	)

	// Validate partially set values
	if req.Input.Prompt != customPrompt {
		t.Errorf("expected Prompt to be %q, got %q", customPrompt, req.Input.Prompt)
	}

	if req.Input.SafetyTolerance != customTolerance {
		t.Errorf("expected SafetyTolerance to be %d, got %d", customTolerance, req.Input.SafetyTolerance)
	}

	// Validate defaults for other fields
	if req.Input.Ratio != DefaultRatio.String() {
		t.Errorf("expected default Ratio to be %q, got %q", DefaultRatio, req.Input.Ratio)
	}

	if req.Input.Raw != DefaultRaw {
		t.Errorf("expected default Raw to be %v, got %v", DefaultRaw, req.Input.Raw)
	}

	if req.Input.ImagePromptStrength != DefaultImagePromptStrength {
		t.Errorf("expected default ImagePromptStrength to be %f, got %f", DefaultImagePromptStrength, req.Input.ImagePromptStrength)
	}
}
