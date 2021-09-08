package validate

import (
	"strings"
	"testing"
)

func TestApplicationName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "h",
			expected: false,
		},
		{
			input:    "he",
			expected: false,
		},
		{
			input:    "hel",
			expected: true,
		},
		{
			input:    "hel2",
			expected: true,
		},
		{
			input:    "_hello",
			expected: false,
		},
		{
			input:    "hello-",
			expected: true,
		},
		{
			input:    "malcolm-in!the-middle",
			expected: false,
		},
		{
			input:    strings.Repeat("a", 63),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 64),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 65),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ApplicationName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
