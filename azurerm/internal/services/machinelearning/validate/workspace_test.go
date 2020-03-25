package validate

import "testing"

func TestWorkspaceName(t *testing.T) {
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
			input:    "hello",
			expected: true,
		},
		{
			// cannot start with a hyphen
			input:    "-hello",
			expected: false,
		},
		{
			// can end with a hyphen
			input:    "hello-",
			expected: true,
		},
		{
			// cannot contain other special symbols other than hyphens
			input:    "hello.world",
			expected: false,
		},
		{
			// hyphen in the middle
			input:    "hello-world",
			expected: true,
		},
		{
			// 2 chars
			input:    "ab",
			expected: false,
		},
		{
			// 3 chars
			input:    "abc",
			expected: true,
		},
		{
			// 33 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefg",
			expected: true,
		},
		{
			// 34 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefgh",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := WorkspaceName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
