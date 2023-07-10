// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = CommunityGalleryImageId{}

func TestCommunityGalleryImageIDFormatter(t *testing.T) {
	actual := NewCommunityGalleryImageID("myGallery1", "myImage1").ID()
	expected := "/communityGalleries/myGallery1/images/myImage1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestCommunityGalleryImageID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CommunityGalleryImageId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing CommunityGalleries
			Input: "/",
			Error: true,
		},

		{
			// missing value for CommunityGalleries
			Input: "/communityGalleries/",
			Error: true,
		},

		{
			// missing images
			Input: "/communityGalleries/myGallery1/",
			Error: true,
		},

		{
			// missing value for images
			Input: "/communityGalleries/myGallery1/images/",
			Error: true,
		},

		{
			// valid
			Input: "/communityGalleries/myGallery1/images/myImage1",
			Expected: &CommunityGalleryImageId{
				GalleryName: "myGallery1",
				ImageName:   "myImage1",
			},
		},

		{
			// upper-cased
			Input: "/COMMUNITYGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CommunityGalleryImageID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.GalleryName != v.Expected.GalleryName {
			t.Fatalf("Expected %q but got %q for GalleryName", v.Expected.GalleryName, actual.GalleryName)
		}

		if actual.ImageName != v.Expected.ImageName {
			t.Fatalf("Expected %q but got %q for ImageName", v.Expected.ImageName, actual.ImageName)
		}
	}
}
