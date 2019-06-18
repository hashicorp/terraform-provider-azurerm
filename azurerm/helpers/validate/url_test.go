package validate

import (
	"testing"
)

func TestURLIsHTTPS(t *testing.T) {
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
			Errors: 1,
		},
		{
			Url:    "https://www.example.com",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Url, func(t *testing.T) {
			_, errors := URLIsHTTPS(tc.Url, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected URLIsHTTPS to have %d not %d errors for %q", tc.Errors, len(errors), tc.Url)
			}
		})
	}
}

func TestURLIsHTTPOrHTTPS(t *testing.T) {
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
		t.Run(tc.Url, func(t *testing.T) {
			_, errors := URLIsHTTPOrHTTPS(tc.Url, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected URLIsHTTPOrHTTPS to have %d not %d errors for %q", tc.Errors, len(errors), tc.Url)
			}
		})
	}
}
