package validate

import "testing"

func TestDataboxEdgePhoneNumber(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "123-555-6789",
			expected: true,
		},
		{
			input:    "+1 123-555-6789",
			expected: true,
		},
		{
			input:    "+92 123 555-6789",
			expected: true,
		},
		{
			input:    "+92 (123) 555-6789",
			expected: true,
		},
		{
			input:    "+92 (123)555-6789",
			expected: true,
		},
		{
			input:    "+1 800-555-6789",
			expected: true,
		},
		{
			input:    "+1-800-555-6789",
			expected: true,
		},
		{
			input:    "1-800-555-6789",
			expected: true,
		},
		{
			input:    "1 800 555-6789",
			expected: true,
		},
		{
			input:    "800-555-6789",
			expected: true,
		},
		{
			input:    "(800)555-6789",
			expected: true,
		},
		{
			input:    "(800) 555-6789",
			expected: true,
		},
		{
			input:    "Jones BBQ and Foot Massage",
			expected: false,
		},
		{
			input:    "800  555-6789",
			expected: false,
		},
		{
			input:    " 1-800-555-6789",
			expected: false,
		},
		{
			input:    "-1-800-555-6789",
			expected: false,
		},
		{
			input:    "123-555 6789",
			expected: false,
		},
		{
			input:    "555-6789",
			expected: false,
		},
		{
			input:    "+12",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgePhoneNumber(v.input, "phone_number")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
