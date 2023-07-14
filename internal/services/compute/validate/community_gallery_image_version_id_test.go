// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestCommunityGalleryImageVersionID(t *testing.T) {
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
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/latest",
			Valid: true,
		},

		{
			// valid
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/1.2.3",
			Valid: true,
		},

		{
			// invalid
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/notTheLatest",
			Valid: false,
		},

		{
			// invalid
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/1.2.",
			Valid: false,
		},

		{
			// invalid
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/1.2",
			Valid: false,
		},

		{
			// invalid
			Input: "/communityGalleries/myGallery1/images/myImage1/versions/1.two.3",
			Valid: false,
		},

		{
			// upper-cased
			Input: "/COMMUNITYGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1/VERSIONS/LATEST",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := CommunityGalleryImageVersionID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t: %+v", tc.Valid, valid, errors)
		}
	}
}
