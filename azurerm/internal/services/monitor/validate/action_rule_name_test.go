package validate

import (
	"testing"
)

func TestActionRuleName(t *testing.T) {
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
			// basic example, lead by lower case letter
			input:    "a",
			expected: true,
		},
		{
			// basic example, lead by upper case letter
			input:    "A",
			expected: true,
		},
		{
			// basic example, lead by number
			input:    "8",
			expected: true,
		},
		{
			// basic example, contain underscore
			input:    "a_b",
			expected: true,
		},
		{
			// basic example, end with underscore
			input:    "ab_",
			expected: true,
		},
		{
			// basic example, end with hyphen
			input:    "ab-",
			expected: true,
		},
		{
			// can not contain '+'
			input:    "a+",
			expected: false,
		},
		{
			// can't lead by '-'
			input:    "-a",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ActionRuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
