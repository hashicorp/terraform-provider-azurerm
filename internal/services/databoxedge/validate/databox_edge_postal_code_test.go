package validate

import "testing"

func TestDataboxEdgePostalCode(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "12345",
			expected: true,
		},
		{
			input:    "12345-1234",
			expected: true,
		},
		{
			input:    "123451234",
			expected: false,
		},
		{
			input:    "12345-123",
			expected: false,
		},
		{
			input:    " 12345-1234",
			expected: false,
		},
		{
			input:    "12345-1234 ",
			expected: false,
		},
		{
			input:    "12345 1234",
			expected: false,
		},
		{
			input:    "123456",
			expected: false,
		},
		{
			input:    "1234",
			expected: false,
		},
		{
			input:    "ManBearPig",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgePostalCode(v.input, "postal_code")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
