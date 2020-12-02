package validate

import "testing"

func TestValidateEmbeddedAdministratorName(t *testing.T) {
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
			input:    "hello",
			expected: false,
		},
		{
			// valid email address
			input:    "hello@microsoft.com",
			expected: true,
		},
		{
			// invalid email address
			input:    "#@%^%#$@#$@#.com",
			expected: false,
		},
		{
			// valid uuid
			input:    "1cf9c591-172b-4654-8ab8-81964aa5335e",
			expected: true,
		},
		{
			// invalid uuid
			input:    "1cf9c591-172b-4654-8ab8-81964aa5335e-0000",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := PowerBIEmbeddedAdministratorName(v.input, "administrators")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
