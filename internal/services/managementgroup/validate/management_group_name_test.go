package validate

import "testing"

func TestManagementGroupName(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "Proper management group name",
			Input:    "Hello",
			Expected: true,
		},
		{
			Name:     "Braces allowed",
			Input:    "Hello(world)",
			Expected: true,
		},
		{
			Name:     "Period allowed",
			Input:    "Hello.world",
			Expected: true,
		},
		{
			Name:     "Hyphen allowed",
			Input:    "Hello-world",
			Expected: true,
		},
		{
			Name:     "Underscore allowed",
			Input:    "Hello_world",
			Expected: true,
		},
		{
			Name:     "Asterisk not allowed",
			Input:    "hello*world",
			Expected: false,
		},
		{
			Name:     "Comma not allowed",
			Input:    "Hello,world",
			Expected: false,
		},
		{
			Name:     "Space not allowed",
			Input:    "Hello world",
			Expected: false,
		},
		{
			Name:     "90 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			Expected: true,
		},
		{
			Name:     "91 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijk",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.Name)

		_, errors := ManagementGroupName(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
