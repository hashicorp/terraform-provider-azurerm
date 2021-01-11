package validate

import (
	"testing"
)

func TestRepoRootFolder(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "/",
			Valid: true,
		},
		{
			Input: "/1",
			Valid: true,
		},
		{
			Input: "/path",
			Valid: true,
		},
		{
			Input: "/path/",
			Valid: true,
		},
		{
			Input: "/path/path",
			Valid: true,
		},
		{
			Input: "/path/path/",
			Valid: true,
		},
		{
			Input: "/path-path",
			Valid: true,
		},
		{
			Input: "/p a t h",
			Valid: true,
		},
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "path",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RepoRootFolder()(tc.Input, "name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
