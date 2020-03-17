package validate

import (
	"testing"
)

func TestCdnEndpointDeliveryPolicyRuleName(t *testing.T) {
	cases := []struct {
		Name        string
		ShouldError bool
	}{
		{
			Name:        "",
			ShouldError: true,
		},
		{
			Name:        "a",
			ShouldError: false,
		},
		{
			Name:        "Z",
			ShouldError: false,
		},
		{
			Name:        "3",
			ShouldError: true,
		},
		{
			Name:        "abc123",
			ShouldError: false,
		},
		{
			Name:        "aBc123",
			ShouldError: false,
		},
		{
			Name:        "aBc 123",
			ShouldError: true,
		},
		{
			Name:        "aBc&123",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := EndpointDeliveryRuleName()(tc.Name, "name")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Name)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Name, len(errors))
			}
		})
	}
}

func TestRuleActionUrlRedirectPath(t *testing.T) {
	cases := []struct {
		Path        string
		ShouldError bool
	}{
		{
			Path:        "",
			ShouldError: false,
		},
		{
			Path:        "a",
			ShouldError: true,
		},
		{
			Path:        "/",
			ShouldError: false,
		},
		{
			Path:        "/abc",
			ShouldError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Path, func(t *testing.T) {
			_, errors := RuleActionUrlRedirectPath()(tc.Path, "path")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Path)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Path, len(errors))
			}
		})
	}
}

func TestRuleActionUrlRedirectQueryString(t *testing.T) {
	cases := []struct {
		QueryString string
		ShouldError bool
	}{
		{
			QueryString: "",
			ShouldError: false,
		},
		{
			QueryString: "a",
			ShouldError: true,
		},
		{
			QueryString: "&a=b",
			ShouldError: true,
		},
		{
			QueryString: "?a=b",
			ShouldError: true,
		},
		{
			QueryString: "a=b",
			ShouldError: false,
		},
		{
			QueryString: "a=b&",
			ShouldError: false,
		},
		{
			QueryString: "a=b&c=d",
			ShouldError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.QueryString, func(t *testing.T) {
			_, errors := RuleActionUrlRedirectQueryString()(tc.QueryString, "query_string")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.QueryString)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.QueryString, len(errors))
			}
		})
	}
}

func TestRuleActionUrlRedirectFragment(t *testing.T) {
	cases := []struct {
		Fragment    string
		ShouldError bool
	}{
		{
			Fragment:    "",
			ShouldError: false,
		},
		{
			Fragment:    "a",
			ShouldError: false,
		},
		{
			Fragment:    "#",
			ShouldError: true,
		},
		{
			Fragment:    "#5fgdfg",
			ShouldError: true,
		},
		{
			Fragment:    "5fgdfg",
			ShouldError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Fragment, func(t *testing.T) {
			_, errors := RuleActionUrlRedirectFragment()(tc.Fragment, "fragment")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Fragment)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Fragment, len(errors))
			}
		})
	}
}

func TestRuleActionCacheExpirationDuration(t *testing.T) {
	cases := []struct {
		Duration    string
		ShouldError bool
	}{
		{
			Duration:    "",
			ShouldError: true,
		},
		{
			Duration:    "23:44:21",
			ShouldError: false,
		},
		{
			Duration:    "9.23:44:21",
			ShouldError: false,
		},
		{
			Duration:    ".23:44:21",
			ShouldError: true,
		},
		{
			Duration:    "9.24:44:21",
			ShouldError: true,
		},
		{
			Duration:    "9.12:61:21",
			ShouldError: true,
		},
		{
			Duration:    "9.11:44:91",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Duration, func(t *testing.T) {
			_, errors := RuleActionCacheExpirationDuration()(tc.Duration, "duration")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Duration)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Duration, len(errors))
			}
		})
	}
}

func TestRuleActionUrlRewriteSourcePattern(t *testing.T) {
	cases := []struct {
		SourcePattern string
		ShouldError   bool
	}{
		{
			SourcePattern: "",
			ShouldError:   true,
		},
		{
			SourcePattern: "a",
			ShouldError:   true,
		},
		{
			SourcePattern: "/",
			ShouldError:   false,
		},
		{
			SourcePattern: "/abc",
			ShouldError:   false,
		},
		{
			SourcePattern: "/abc\n",
			ShouldError:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.SourcePattern, func(t *testing.T) {
			_, errors := RuleActionUrlRewriteSourcePattern()(tc.SourcePattern, "source_pattern")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.SourcePattern)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.SourcePattern, len(errors))
			}
		})
	}
}

func TestRuleActionUrlRewriteSourceDestination(t *testing.T) {
	cases := []struct {
		Destination string
		ShouldError bool
	}{
		{
			Destination: "",
			ShouldError: true,
		},
		{
			Destination: "a",
			ShouldError: true,
		},
		{
			Destination: "/",
			ShouldError: false,
		},
		{
			Destination: "/abc",
			ShouldError: false,
		},
		{
			Destination: "/abc\n",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Destination, func(t *testing.T) {
			_, errors := RuleActionUrlRewriteDestination()(tc.Destination, "destination")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Destination)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Destination, len(errors))
			}
		})
	}
}
