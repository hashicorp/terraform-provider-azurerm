package validate

import "testing"

func TestBlueprintDefinitionName(t *testing.T) {
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
			Name:     "Basic example",
			Input:    "hello",
			Expected: true,
		},
		{
			Name:     "Cannot contain underscore",
			Input:    "_hello",
			Expected: false,
		},
		{
			Name:     "Cannot contain dot",
			Input:    ".hello",
			Expected: false,
		},
		{
			Name:     "Start with hyphen",
			Input:    "-helloworld",
			Expected: true,
		},
		{
			Name:     "Hyphen in middle",
			Input:    "hello-world",
			Expected: true,
		},
		{
			Name:     "End with hyphen",
			Input:    "helloworld-",
			Expected: true,
		},
		{
			Name:     "Cannot contain an exclamation mark",
			Input:    "hello!",
			Expected: false,
		},
		{
			Name:     "48 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefgh",
			Expected: true,
		},
		{
			Name:     "49 characters",
			Input:    "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghi",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.Name)

		_, errors := BlueprintDefinitionName(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
