package validate

import (
	"testing"
)

func TestApplicationGatewayName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "a",
			expected: true,
		},
		{
			input:    "1",
			expected: true,
		},
		{
			input:    "_",
			expected: false,
		},
		{
			input:    "a_",
			expected: true,
		},
		{
			input:    "a-._0",
			expected: true,
		},
		{
			input:    "-abc",
			expected: false,
		},
		{
			input:    "abc-",
			expected: false,
		},
		{
			input:    "abc!",
			expected: false,
		},
		{
			input:    "valid_name.123-abc",
			expected: true,
		},
		{
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab",
			expected: true,
		},
		{
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc",
			expected: false,
		},
	}

	for _, tc := range tests {
		_, errs := ApplicationGatewayName(tc.input, "name")
		ok := len(errs) == 0
		if ok != tc.expected {
			t.Fatalf("input %q: expected %t got %t", tc.input, tc.expected, ok)
		}
	}
}
