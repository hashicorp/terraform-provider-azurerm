package validate

import "testing"

func TestSpringCloudAppName(t *testing.T) {
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
			// can't end with hyphen
			input:    "abc-",
			expected: false,
		},
		{
			// can not short than 4 characters
			input:    "abc",
			expected: false,
		},
		{
			// 32 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdef",
			expected: true,
		},
		{
			// 33 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefg",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := SpringCloudAppName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
