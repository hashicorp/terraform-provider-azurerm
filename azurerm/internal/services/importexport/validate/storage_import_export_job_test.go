package validate

import "testing"

func TestImportExportJobName(t *testing.T) {
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
			// can not short than 2 characters
			input:    "a",
			expected: false,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			expected: true,
		},
		{
			// 65 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ImportExportJobName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestImportExportJobPhone(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "hello123",
			expected: false,
		},
		{
			input:    "hello!",
			expected: false,
		},
		{
			input:    "123hello123",
			expected: false,
		},
		{
			input:    "123-242345",
			expected: true,
		},
		{
			input:    "+12",
			expected: true,
		},
		{
			input:    "+12 123456789",
			expected: true,
		},
		{
			input:    "1123456789",
			expected: true,
		},
		{
			input:    "+(1) 23456789",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ImportExportJobPhone(v.input, "phone")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
