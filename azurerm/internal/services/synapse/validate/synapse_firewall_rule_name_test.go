package validate

import (
	"testing"
)

func TestSynapseFirewallRuleName(t *testing.T) {
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
			// can contain underscore
			input:    "aBc_123",
			expected: true,
		},
		{
			// can contain hyphen
			input:    "ab-c",
			expected: true,
		},
		{
			// can't contain `*`
			input:    "abcon*demand",
			expected: false,
		},
		{
			// 128 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdef",
			expected: true,
		},
		{
			// 129 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefg",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := SynapseFirewallRuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
