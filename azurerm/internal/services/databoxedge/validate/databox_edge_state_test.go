package validate

import "testing"

func TestDataboxEdgeState(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "W",
			expected: true,
		},
		{
			input:    "Washington",
			expected: true,
		},
		{
			input:    "Pekwachnamaykoskwaskwaypinwanik",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeState(v.input, "state")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
