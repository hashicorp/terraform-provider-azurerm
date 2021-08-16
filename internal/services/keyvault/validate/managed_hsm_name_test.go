package validate

import "testing"

func TestManagedHardwareSecurityModuleName(t *testing.T) {
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
			input:    "abcd",
			expected: true,
		},
		{
			// can't start with a number
			input:    "1abc",
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
			// can not short than 3 characters
			input:    "ab",
			expected: false,
		},
		{
			// 24 chars
			input:    "abcdefghijklmnopqrstuvwx",
			expected: true,
		},
		{
			// 25 chars
			input:    "abcdefghijklmnopqrstuvwxy",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ManagedHardwareSecurityModuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
