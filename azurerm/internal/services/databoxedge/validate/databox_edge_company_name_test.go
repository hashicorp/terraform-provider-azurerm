package validate

import "testing"

func TestDataboxEdgeCompanyName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "A",
			expected: false,
		},
		{
			input:    "Mi",
			expected: true,
		},
		{
			input:    "Acme Corporation",
			expected: true,
		},
		{
			input:    "_Stark Industries",
			expected: true,
		},
		{
			input:    "Wayne Enterprises-",
			expected: true,
		},
		{
			input:    "Flynns Arcade!",
			expected: true,
		},
		{
			input:    "Ollivander's Wand Shop",
			expected: true,
		},
		{
			input:    "Cyberdyne Systems.",
			expected: true,
		},
		{
			input:    "Genco Pura Olive Oil Company",
			expected: true,
		},
		{
			input:    "This is loongest valid Company Name",
			expected: true,
		},
		{
			input:    "This is Company Name is just to long",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeCompanyName(v.input, "company_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
