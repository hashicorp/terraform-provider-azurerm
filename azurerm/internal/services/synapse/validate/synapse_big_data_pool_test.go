package validate

import (
	"testing"
)

func TestSynapseBigDataPoolName(t *testing.T) {
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
			input:    "aBc123",
			expected: true,
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
			// can't start with a number
			input:    "123abc",
			expected: false,
		},
		{
			// 15 chars
			input:    "abcdefghijklmno",
			expected: true,
		},
		{
			// 16 chars
			input:    "abcdefghijklmnop",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := SynapseBigDataPoolName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
