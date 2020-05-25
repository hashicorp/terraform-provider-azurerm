package azure

import (
	"strings"
	"testing"
)

func TestDataBoxJobName(t *testing.T) {
	testData := []struct {
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
			input:    "_a",
			expected: false,
		},
		{
			input:    "a-",
			expected: false,
		},
		{
			input:    "12345",
			expected: true,
		},
		{
			input:    "1",
			expected: true,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			input:    strings.Repeat("w", 259),
			expected: true,
		},
		{
			input:    strings.Repeat("w", 260),
			expected: true,
		},
		{
			input:    strings.Repeat("w", 261),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateServiceBusTopicName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
