package validate

import "testing"

func TestSpringCloudCustomDomainName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "golang.org",
			Valid: true,
		},
		{
			Input: "golang-ci.org",
			Valid: true,
		},
		{
			Input: "golang.vip",
			Valid: true,
		},
		{
			Input: "golang.tokyo",
			Valid: true,
		},
		{
			Input: "www.golang.org",
			Valid: true,
		},
		{
			Input: "subdomain.golang.org",
			Valid: true,
		},
		{
			Input: "www.golang.vip",
			Valid: true,
		},
		{
			Input: "www.golang.tokyo",
			Valid: true,
		},
		{
			Input: "host.subdomain.golang.org",
			Valid: true,
		},
		{
			Input: ".golang-ci.org",
			Valid: false,
		},
		{
			Input: "golang",
			Valid: false,
		},
		{
			Input: "golang.org.",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SpringCloudCustomDomainName(tc.Input, "name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
