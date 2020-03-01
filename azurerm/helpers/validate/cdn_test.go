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
			_, errors := CdnEndpointDeliveryPolicyRuleName()(tc.Name, "name")

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
