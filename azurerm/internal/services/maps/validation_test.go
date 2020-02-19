package maps

import "testing"

func TestValidateName(t *testing.T) {
	testData := []struct {
		Name     string
		Expected bool
	}{
		{
			Name:     "",
			Expected: false,
		},
		{
			Name:     "hello",
			Expected: true,
		},
		{
			Name:     "Hello",
			Expected: true,
		},
		{
			Name:     "1hello",
			Expected: true,
		},
		{
			Name:     "1he-llo",
			Expected: true,
		},
		{
			Name:     "he-llo1",
			Expected: true,
		},
		{
			Name:     "he_llo1",
			Expected: true,
		},
		{
			Name:     ".hello1",
			Expected: false,
		},
		{
			Name:     "_hello1",
			Expected: false,
		},
		{
			Name:     "he.llo1",
			Expected: true,
		},
		{
			Name:     "he-llo!",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		warnings, errors := ValidateName()(v.Name, "name")
		if len(warnings) != 0 {
			t.Fatalf("Expected no warnings but got %d", len(warnings))
		}

		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t for %q: %s", v.Expected, actual, v.Name, errors)
		}
	}
}
