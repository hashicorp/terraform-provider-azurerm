package validate

import (
	"testing"
)

func TestPrivateLinkHubName(t *testing.T) {
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
			input:    "abc123",
			expected: true,
		},
		{
			// can't contain upper case
			input:    "aBc123",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "abc_123",
			expected: false,
		},
		{
			// can't contain hyphen
			input:    "abc-123",
			expected: false,
		},
		{
			// 45 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrs",
			expected: true,
		},
		{
			// 46 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := PrivateLinkHubName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
