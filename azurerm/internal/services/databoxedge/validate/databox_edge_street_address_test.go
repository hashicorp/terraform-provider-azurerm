package validate

import "testing"

func TestDataboxEdgeStreetAddress(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "12 Grimmauld Place",
			expected: true,
		},
		{
			input:    "740 Evergreen Terrace",
			expected: true,
		},
		{
			input:    "129 West 81st Street , Apartment: 5A",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeStreetAddress(v.input, "street_address")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
