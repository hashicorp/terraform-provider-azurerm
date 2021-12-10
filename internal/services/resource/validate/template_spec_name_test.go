package validate

import "testing"

func TestTemplateSpecName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			// invalid char
			Input: "{scwiffy}",
			Valid: false,
		},
		{
			// too long - 65 chars
			Input: "01234567890123456789012345678901234567890123456789012345678901234",
			Valid: false,
		},
		{
			// short alpha
			Input: "a",
			Valid: true,
		},
		{
			// valid special
			Input: "(",
			Valid: true,
		},
		{
			// sensible value
			Input: "Production",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing value %s", tc.Input)
		_, errors := TemplateSpecName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("expected %t but got %t", tc.Valid, valid)
		}
	}
}
