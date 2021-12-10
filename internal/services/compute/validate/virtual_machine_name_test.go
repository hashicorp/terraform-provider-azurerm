package validate

import "testing"

func TestVirtualMachineName(t *testing.T) {
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
			// 79 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyza",
			expected: true,
		},
		{
			// 80 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab",
			expected: true,
		},
		{
			// 81 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc",
			expected: false,
		},
		{
			// may contain alphanumerics, dots, dashes and underscores
			input:    "hello_world7.goodbye-world4",
			expected: true,
		},
		{
			// must begin with an alphanumeric
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can end with an underscore
			input:    "hello_",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// start with a number
			input:    "0abc",
			expected: true,
		},
		{
			// cannot contain only numbers
			input:    "12345",
			expected: false,
		},
		{
			// can start with upper case letter
			input:    "Test",
			expected: true,
		},
		{
			// can end with upper case letter
			input:    "TEST",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := VirtualMachineName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
