package validate

import "testing"

func TestWindowsComputerName(t *testing.T) {
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
			// can't contain underscore
			input:    "hello_world",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// can't contain dot
			input:    "hello.world",
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
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := windowsComputerName(v.input, "computer_name", 100)
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestWindowsComputerNameFull(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// 14 chars
			input:    "abcdefghijklmn",
			expected: true,
		},
		{
			// 15 chars
			input:    "abcdefghijklmno",
			expected: true,
		},
		{
			// 16 chars
			input:    "abcdefghijklmnop",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := WindowsComputerNameFull(v.input, "computer_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestWindowsComputerNamePrefix(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// 8 chars
			input:    "abcdefgh",
			expected: true,
		},
		{
			// 9 chars
			input:    "abcdefghi",
			expected: true,
		},
		{
			// 10 chars
			input:    "abcdefghij",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := WindowsComputerNamePrefix(v.input, "computer_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
