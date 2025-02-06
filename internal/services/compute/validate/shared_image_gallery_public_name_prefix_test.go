// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestSharedImageGalleryPrefix(t *testing.T) {
	validNames := []string{
		"Valid012",
		"abcde",
		"0123456789012345",
	}
	for _, v := range validNames {
		_, errors := SharedImageGalleryPrefix(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Shared Image Gallery Public Name Prefix: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"-",
		"abcd",
		"01234567890123456",
	}
	for _, v := range invalidNames {
		_, errors := SharedImageGalleryPrefix(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Shared Image Gallery Public Name Prefix", v)
		}
	}
}
