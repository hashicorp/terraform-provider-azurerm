package validate

import "testing"

func TestAutomationConnectionName(t *testing.T) {
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
			// contain underscore and hyphens
			input:    "ab_c-abc",
			expected: true,
		},
		{
			// starts with underscore
			input:    "_abc",
			expected: true,
		},
		{
			// starts with hyphens
			input:    "-abc",
			expected: true,
		},
		{
			// starts with number
			input:    "1abc",
			expected: true,
		},
		{
			// can't contain +
			input:    "ab+c",
			expected: false,
		},
		{
			// can't end with %
			input:    "ab%c",
			expected: false,
		},
		{
			// 128 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdef",
			expected: true,
		},
		{
			// 129 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdefa",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ConnectionName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
