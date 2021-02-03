package validate

import (
	"testing"
)

func TestSharedAccessSignatureIP(t *testing.T) {
	cases := []struct {
		Input       string
		ShouldError bool
	}{
		{
			Input:       "dfdsdfds",
			ShouldError: true,
		},
		{
			Input:       "192.168.0.1",
			ShouldError: false,
		},
		{
			Input:       "sdfsdf-4334",
			ShouldError: true,
		},
		{
			Input:       "172.77.62-abc",
			ShouldError: true,
		},
		{
			Input:       "66.247.118.148-66.247.118.148",
			ShouldError: true,
		},
		{
			Input:       "66.247.118.148-66.247.118.149",
			ShouldError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := SharedAccessSignatureIP(tc.Input, "ip")

			hasErrors := len(errors) > 0
			if !hasErrors && tc.ShouldError {
				t.Fatalf("Expected an error but didn't get one for %q", tc.Input)
			}

			if hasErrors && !tc.ShouldError {
				t.Fatalf("Expected to get no errors for %q but got %d", tc.Input, len(errors))
			}
		})
	}
}
