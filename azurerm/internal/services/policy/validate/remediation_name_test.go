package validate

import "testing"

func TestValidateRemediationName(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "empty",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "basic example",
			Input:    "hello",
			Expected: true,
		},
		{
			Name:     "start with an underscore",
			Input:    "_hello",
			Expected: true,
		},
		{
			Name:     "end with a hyphen",
			Input:    "hello-",
			Expected: true,
		},
		{
			Name:     "can contain an exclamation mark",
			Input:    "hello!",
			Expected: true,
		},
		{
			Name:     "dash in the middle",
			Input:    "malcolm-middle",
			Expected: true,
		},
		{
			Name:     "can't end with a period",
			Input:    "hello.",
			Expected: true,
		},
		{
			Name:     "cannot contain %",
			Input:    "hello%world",
			Expected: false,
		},
		{
			Name:     "cannot contain ^",
			Input:    "hello^world",
			Expected: false,
		},
		{
			Name:     "cannot contain #",
			Input:    "hello#world",
			Expected: false,
		},
		{
			Name:     "cannot contain ?",
			Input:    "hello?world",
			Expected: false,
		},
		{
			Name:     "cannot contains upper case letters",
			Input:    "HelloWorld",
			Expected: false,
		},
		{
			Name:     "260 chars",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			Expected: true,
		},
		{
			Name:     "261 chars",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijk",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		_, errors := RemediationName(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
