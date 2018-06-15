package validate

import "testing"

func TestHelper_Validate_StringNotEmpty(t *testing.T) {
	cases := []struct {
		String string
		Errors int
	}{
		{
			String: "",
			Errors: 1,
		},
		{
			String: "k",
			Errors: 0,
		},
		{
			String: "a longer string",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := StringNotEmpty(tc.String, "test")

		if len(errors) < tc.Errors {
			t.Fatalf("Expected StringNotEmpty to have an error for '%q'", tc.String)
		}
	}
}
