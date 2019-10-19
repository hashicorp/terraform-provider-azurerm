package validate

import "testing"

func TestValidateRFC3986Attribute(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "_",
			Expected: false,
		},
		{
			Input:    "?#$~()*&^%[]/<>",
			Expected: false,
		},
		{
			Input:    "a",
			Expected: true,
		},
		{
			Input:    "A",
			Expected: true,
		},
		{
			Input:    "123-abc",
			Expected: true,
		},
		{
			Input:    "a-",
			Expected: false,
		},
		{
			Input:    "-",
			Expected: false,
		},
		{
			Input:    "a.-_",
			Expected: true,
		},
		{
			Input:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890",
			Expected: true,
		},
		{
			Input:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890_",
			Expected: false,
		},
		{
			Input:    "aA_-.1",
			Expected: true,
		},
		{
			Input:    "a _",
			Expected: false,
		},
	}

	for _, v := range testCases {
		t.Logf("[DEBUG] Test Input %q", v.Input)

		warnings, errors := GenericRFC3986Compliance(v.Input, "RFC3986Attribute")
		if len(warnings) != 0 {
			t.Fatalf("Expected no warnings but got %d", len(warnings))
		}

		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
