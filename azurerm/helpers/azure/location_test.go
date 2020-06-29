package azure

import "testing"

func TestNormalizeLocation(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "West US",
			expected: "westus",
		},
		{
			input:    "South East Asia",
			expected: "southeastasia",
		},
		{
			input:    "southeastasia",
			expected: "southeastasia",
		},
	}

	for _, v := range cases {
		actual := NormalizeLocation(v.input)
		if v.expected != actual {
			t.Fatalf("Expected %q but got %q", v.expected, actual)
		}
	}
}
