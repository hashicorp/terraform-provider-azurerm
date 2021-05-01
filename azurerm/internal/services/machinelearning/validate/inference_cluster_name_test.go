package validate

import "testing"

func TestInferenceClusterName(t *testing.T) {
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
			// 1 char
			input:    "a",
			expected: false,
		},
		{
			// 2 chars
			input:    "ab",
			expected: true,
		},
		{
			// 16 chars
			input:    "abcdefghijklmnop",
			expected: true,
		},
		{
			// 17 chars
			input:    "abcdefghijklmnopq",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := InferenceClusterName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
