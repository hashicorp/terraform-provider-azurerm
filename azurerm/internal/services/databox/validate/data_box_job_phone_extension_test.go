package validate

import "testing"

func TestDataBoxJobPhoneExtension(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: true,
		},
		{
			input:    "hello",
			expected: false,
		},
		{
			input:    "123",
			expected: true,
		},
		{
			input:    "1234",
			expected: true,
		},
		{
			input:    "12345",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataBoxJobPhoneExtension(v.input, "phone_extension")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
