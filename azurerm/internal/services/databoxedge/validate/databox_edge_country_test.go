package validate

import "testing"

func TestDataboxEdgeCountry(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "c",
			expected: false,
		},
		{
			input:    "United Kingdom",
			expected: true,
		},
		{
			input:    "US",
			expected: false,
		},
		{
			input:    "UsA",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeCountry(v.input, "country")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
