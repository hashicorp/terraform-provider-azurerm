package validate

import "testing"

func TestUrlTemplate(t *testing.T) {
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
			// not a url
			Input: "not a url",
			Valid: false,
		},

		{
			// contains http
			Input: "http",
			Valid: false,
		},

		{
			// not correct
			Input: "http:/",
			Valid: false,
		},

		{
			// empty
			Input: "http://",
			Valid: false,
		},

		{
			// right
			Input: "http://abc.com",
			Valid: true,
		},

		{
			// multi path
			Input: "https://abc.com/api/test",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)

		_, errors := UrlTemplate(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
