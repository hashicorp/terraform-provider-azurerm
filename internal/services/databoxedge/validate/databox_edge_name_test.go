package validate

import "testing"

func TestDataboxEdgeName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "hello",
			expected: true,
		},
		{
			input:    "_hello",
			expected: false,
		},
		{
			input:    "hello-",
			expected: false,
		},
		{
			input:    "hello!",
			expected: false,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			input:    "hello.",
			expected: false,
		},
		{
			input:    "qwertyuioplkjhgfdsazxcv",
			expected: true,
		},
		{
			input:    "qwertyuioplkjhgfdsazxcva",
			expected: true,
		},
		{
			input:    "qwertyuioplkjhgfdsazxcvgg",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
