package validate

import (
	"testing"
)

func TestWorkspaceName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// 1 char
			input:    "a",
			expected: true,
		},
		{
			// basic example
			input:    "abc-123",
			expected: true,
		},
		{
			// can't contain upper case
			input:    "aBc123",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "ab_c",
			expected: false,
		},
		{
			// can't contain `-ondemand`
			input:    "abc-ondemand123",
			expected: false,
		},
		{
			// must start lowercase char or number
			input:    "-abc123",
			expected: false,
		},
		{
			// must end lowercase char or number
			input:    "abc123-",
			expected: false,
		},
		{
			// 50 chars
			input:    "abcdefghijklmnopqrstuvwxyz-abcdefghijklmnopqrstuvw",
			expected: true,
		},
		{
			// 51 chars
			input:    "abcdefghijklmnopqrstuvwxyz-abcdefghijklmnopqrstuvwx",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := WorkspaceName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
