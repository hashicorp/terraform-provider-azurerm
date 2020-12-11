package validate

import "testing"

func TestDataBoxJobStreetAddress(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "16 TOWNSEND ST",
			expected: true,
		},
		{
			input:    "qwertyuiopasdasgfdhdghjjkljklzxcxbc",
			expected: true,
		},
		{
			input:    "qwertyuiopasdasgfdhdghjjkljklzxcxbcz",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataBoxJobStreetAddress(v.input, "street_address")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
