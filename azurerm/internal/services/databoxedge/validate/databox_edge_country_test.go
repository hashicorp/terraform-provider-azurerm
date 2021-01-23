package validate

import "testing"

func TestDataboxEdgeCounty(t *testing.T) {
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
			expected: false,
		},
		{
			input:    "US",
			expected: true,
		},
		{
			input:    "UsA",
			expected: true,
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
