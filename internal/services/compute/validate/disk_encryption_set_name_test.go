package validate

import "testing"

func TestDiskEncryptionSetName(t *testing.T) {
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
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't start with dot
			input:    ".hello",
			expected: false,
		},
		{
			// dot in middle
			input:    "hello.world",
			expected: true,
		},
		{
			// hyphen in middle
			input:    "hello-world",
			expected: true,
		},
		{
			// can't end with hyphen
			input:    "helloworld-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// can't end with dot
			input:    "hello.",
			expected: false,
		},
		{
			// underscore at end
			input:    "helloworld_",
			expected: true,
		},
		{
			// 80 characters
			input:    "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			expected: true,
		},
		{
			// 81 characters
			input:    "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdef",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.input)

		_, errors := DiskEncryptionSetName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
