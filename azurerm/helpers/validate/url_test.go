package validate

import (
	"testing"
)

func TestHelper_Validate_UrlIsHttpOrHttps(t *testing.T) {
	cases := []struct {
		Url    string
		Errors int
	}{
		{
			Url:    "",
			Errors: 1,
		},
		{
			Url:    "this is not a url",
			Errors: 1,
		},
		{
			Url:    "www.example.com",
			Errors: 1,
		},
		{
			Url:    "ftp://www.example.com",
			Errors: 1,
		},
		{
			Url:    "http://www.example.com",
			Errors: 0,
		},
		{
			Url:    "https://www.example.com",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := UrlIsHttpOrHttps(tc.Url, "test")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected UrlIsHttpOrHttps to have an error for %q", tc.Url)
		}
	}
}
