package validate

import "testing"

func TestDataBoxJobPhoneNumber(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "+1123456789",
			expected: true,
		},
		{
			input:    "123456789",
			expected: false,
		},
		{
			input:    "hello123",
			expected: false,
		},
		{
			input:    "hello!",
			expected: false,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: false,
		},
		{
			input:    "hello.",
			expected: false,
		},
		{
			input:    "+1",
			expected: false,
		},
		{
			input:    "+12",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataBoxJobPhoneNumber(v.input, "phone_number")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
