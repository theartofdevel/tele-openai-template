package replicate

import (
	"testing"
)

func TestAspectRatio_String(t *testing.T) {
	tests := []struct {
		input    AspectRatio
		expected string
	}{
		{AspectRatio21x9, "21:9"},
		{AspectRatio16x9, "16:9"},
		{AspectRatio3x2, "3:2"},
		{AspectRatio4x3, "4:3"},
		{AspectRatio5x4, "5:4"},
		{AspectRatio1x1, "1:1"},
		{AspectRatio4x5, "4:5"},
		{AspectRatio3x4, "3:4"},
		{AspectRatio2x3, "2:3"},
		{AspectRatio9x16, "9:16"},
		{AspectRatio9x21, "9:21"},
		{-1, "unknown"},               // Invalid input
		{AspectRatio(100), "unknown"}, // Out-of-range input
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		input          string
		expected       AspectRatio
		expectedErrMsg string
	}{
		{"21:9", AspectRatio21x9, ""},
		{"16:9", AspectRatio16x9, ""},
		{"3:2", AspectRatio3x2, ""},
		{"4:3", AspectRatio4x3, ""},
		{"5:4", AspectRatio5x4, ""},
		{"1:1", AspectRatio1x1, ""},
		{"4:5", AspectRatio4x5, ""},
		{"3:4", AspectRatio3x4, ""},
		{"2:3", AspectRatio2x3, ""},
		{"9:16", AspectRatio9x16, ""},
		{"9:21", AspectRatio9x21, ""},
		{"invalid", -1, "invalid aspect ratio string"}, // Invalid string
		{"", -1, "invalid aspect ratio string"},        // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := NewAspectRatio(tt.input)
			if tt.expectedErrMsg == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			} else {
				if err == nil || err.Error() != tt.expectedErrMsg {
					t.Errorf("expected error %q, got %v", tt.expectedErrMsg, err)
				}
			}
		})
	}
}

func TestBidirectionalMapping(t *testing.T) {
	for ar, str := range mapping {
		t.Run(str, func(t *testing.T) {
			// Test String method
			if ar.String() != str {
				t.Errorf("String() for %v expected %q, got %q", ar, str, ar.String())
			}
			// Test FromString
			res, err := NewAspectRatio(str)
			if err != nil {
				t.Errorf("FromString(%q) returned error: %v", str, err)
			}
			if res != ar {
				t.Errorf("FromString(%q) expected %v, got %v", str, ar, res)
			}
		})
	}
}
