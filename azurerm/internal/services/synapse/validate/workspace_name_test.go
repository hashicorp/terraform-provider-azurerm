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
			input:    "ab_c",
			expected: false,
		},
		{
			// can't contain hyphen
			input:    "ab-c",
			expected: false,
		},
		{
			// can't end with `ondemand`
			input:    "abcondemand",
			expected: false,
		},
		{
			// 45 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklm",
			expected: true,
		},
		{
			// 46 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmn",
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
