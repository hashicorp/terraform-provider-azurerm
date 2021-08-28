package validate

import "testing"

func TestIotSecuritySolutionName(t *testing.T) {
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
			// can't contain character other than letter, digit, '-', '.' and '_'
			input:    "ab*",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := IotSecuritySolutionName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
