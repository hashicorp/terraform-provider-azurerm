package validate

import (
	"testing"
)

func TestSharedImageGalleryName(t *testing.T) {
	cases := []struct {
		Input       string
		ShouldError bool
	}{
		{
			Input:       "",
			ShouldError: true,
		},
		{
			Input:       "a.b.c",
			ShouldError: true,
		},
		{
			Input:       "1.2.3",
			ShouldError: false,
		},
		{
			Input:       "0.0.1",
			ShouldError: false,
		},
		{
			Input:       "hello",
			ShouldError: true,
		},
		{
			Input:       "1.2.3.4",
			ShouldError: true,
		},
		{
			Input:       "hell0-there",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := SharedImageVersionName(tc.Input, "test")

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
