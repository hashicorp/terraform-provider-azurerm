package validate

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
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
			Input:       "hello",
			ShouldError: false,
		},
		{
			Input:       "hello123",
			ShouldError: false,
		},
		{
			Input:       "hello.123",
			ShouldError: false,
		},
		{
			Input:       "hello,123",
			ShouldError: true,
		},
		{
			Input:       "hello_123",
			ShouldError: false,
		},
		{
			Input:       "hello-123",
			ShouldError: true,
		},
		{
			Input:       acctest.RandString(79),
			ShouldError: false,
		},
		{
			Input:       acctest.RandString(80),
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := SharedImageGalleryName(tc.Input, "test")

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

func TestSharedImageName(t *testing.T) {
	cases := []struct {
		Input       string
		ShouldError bool
	}{
		{
			Input:       "",
			ShouldError: true,
		},
		{
			Input:       "hello",
			ShouldError: false,
		},
		{
			Input:       "hello123",
			ShouldError: false,
		},
		{
			Input:       "hello.123",
			ShouldError: false,
		},
		{
			Input:       "hello,123",
			ShouldError: true,
		},
		{
			Input:       "hello_123",
			ShouldError: false,
		},
		{
			Input:       "hello-123",
			ShouldError: false,
		},
		{
			Input:       acctest.RandString(79),
			ShouldError: false,
		},
		{
			Input:       acctest.RandString(80),
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := SharedImageName(tc.Input, "test")

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

func TestSharedImageVersionName(t *testing.T) {
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

func TestVirtualMachineTimeZone(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 0,
		},
		{
			Value:  "UTC",
			Errors: 0,
		},
		{
			Value:  "China Standard Time",
			Errors: 0,
		},
		{
			// Valid UTC time zone
			Value:  "utc-11",
			Errors: 0,
		},
		{
			// Invalid UTC time zone
			Value:  "UTC-30",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualMachineTimeZone()(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected VirtualMachineTimeZone to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
