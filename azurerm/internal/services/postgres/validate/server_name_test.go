package validate

import (
	"testing"
)

func TestValidateServerName(t *testing.T) {
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
			input:    "ab-c",
			expected: true,
		},
		{
			// can't contain upper case letter
			input:    "AbcD",
			expected: false,
		},
		{
			// can't start with a hyphen
			input:    "-abc",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "ab_c",
			expected: false,
		},
		{
			// can't end with hyphen
			input:    "abc-",
			expected: false,
		},
		{
			// can not short than 3 characters
			input:    "ab",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcde",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdef",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := PostgreSQLServerName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
