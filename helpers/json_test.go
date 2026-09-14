package helpers

import (
	"strings"
	"testing"
)

func TestNormalizeJson(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "valid JSON (compact)",
			input:    `{"key":"value"}`,
			expected: `{"key":"value"}`,
		},
		{
			name:     "valid JSON (formatted)",
			input:    "{\n  \"key\": \"value\"\n}",
			expected: `{"key":"value"}`,
		},
		{
			name:     "valid JSON (arrays and nested)",
			input:    `{ "a": [1, 2, 3], "b": { "c": "d" } }`,
			expected: `{"a":[1,2,3],"b":{"c":"d"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NormalizeJson(tc.input)
			if actual != tc.expected {
				t.Fatalf("expected: %q, got: %q", tc.expected, actual)
			}
		})
	}
}

func TestNormalizeJson_InvalidJson(t *testing.T) {
	input := `{"key": "value"` // Missing closing brace
	actual := NormalizeJson(input)
	if !strings.HasPrefix(actual, "Error parsing JSON:") {
		t.Fatalf("expected error prefix, got: %q", actual)
	}
}

func TestNormalizeJson_NonString(t *testing.T) {
	actual := NormalizeJson(123)
	if !strings.HasPrefix(actual, "Error parsing JSON: expected string, got int") {
		t.Fatalf("expected error string for non-string input, got: %q", actual)
	}
}
