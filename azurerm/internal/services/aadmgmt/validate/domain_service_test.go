package validate

import "testing"

func TestValidateDomainServiceName(t *testing.T) {
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
			input:    "test.onmicrosoft.com",
			expected: true,
		},
		{
			// can't larger than 64 characters in length
			input:    "abcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyz.abcdefghijk",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "ab_c.com",
			expected: false,
		},
		{
			// can't have only one segment
			input:    "abcd",
			expected: false,
		},
		{
			// the first segment can not be all numbers
			input:    "123.test.com",
			expected: false,
		},
		{
			// each segment can only start with number or letters
			input:    "-a.test.com",
			expected: false,
		},
		{
			// The prefix can not contain more than 15 characters
			input:    "abcdefghijklmnop.test.com",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateDomainServiceName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
