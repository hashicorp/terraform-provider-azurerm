package validate

import (
	"strings"
	"testing"
)

func TestClusterName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "test",
			expected: true,
		},
		{
			input:    "test-abc",
			expected: true,
		},
		{
			input:    strings.Repeat("s", 259),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 260),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 261),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ClusterName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
