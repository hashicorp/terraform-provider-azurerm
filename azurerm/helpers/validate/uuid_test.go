package validate

import "testing"

func TestUUID(t *testing.T) {
	cases := []struct {
		Input  string
		Errors int
	}{
		{
			Input:  "",
			Errors: 1,
		},
		{
			Input:  "hello-world",
			Errors: 1,
		},
		{
			Input:  "00000000-0000-111-0000-000000000000",
			Errors: 1,
		},
		{
			Input:  "00000000-0000-0000-0000-000000000000",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := UUID(tc.Input, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected UUID to have %d not %d errors for %q", tc.Errors, len(errors), tc.Input)
			}
		})
	}
}
