// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestCommunityGalleryImageID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing CommunityGalleries
			Input: "/",
			Valid: false,
		},

		{
			// missing value for CommunityGalleries
			Input: "/communityGalleries/",
			Valid: false,
		},

		{
			// missing images
			Input: "/communityGalleries/myGallery1/",
			Valid: false,
		},

		{
			// missing value for images
			Input: "/communityGalleries/myGallery1/images/",
			Valid: false,
		},

		{
			// valid
			Input: "/communityGalleries/myGallery1/images/myImage1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/COMMUNITYGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CommunityGalleryImageID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t: %+v", tc.Valid, valid, errors)
		}
	}
}
