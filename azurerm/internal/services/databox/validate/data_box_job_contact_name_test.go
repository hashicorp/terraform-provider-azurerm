package validate

import "testing"

func TestDataBoxJobContactName(t *testing.T) {
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
			expected: true,
		},
		{
			input:    "hello-",
			expected: true,
		},
		{
			input:    "hello!",
			expected: true,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			input:    "hello.",
			expected: true,
		},
		{
			input:    "ahellodfasdfsdafsdfsdfsdfasdfsdfb",
			expected: true,
		},
		{
			input:    "ahellodfasdfsdafsdfsdfsdfasdfsdfbc",
			expected: true,
		},
		{
			input:    "ahellodfasdfsdafsdfsdfsdfasdfsdfbss",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataBoxJobContactName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
