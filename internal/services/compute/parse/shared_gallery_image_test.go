// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SharedGalleryImageId{}

func TestSharedGalleryImageIDFormatter(t *testing.T) {
	actual := NewSharedGalleryImageID("myGallery1", "myImage1").ID()
	expected := "/sharedGalleries/myGallery1/images/myImage1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSharedGalleryImageID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SharedGalleryImageId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SharedGalleries
			Input: "/",
			Error: true,
		},

		{
			// missing value for SharedGalleries
			Input: "/sharedGalleries/",
			Error: true,
		},

		{
			// missing images
			Input: "/sharedGalleries/myGallery1/",
			Error: true,
		},

		{
			// missing value for images
			Input: "/sharedGalleries/myGallery1/images/",
			Error: true,
		},

		{
			// valid
			Input: "/sharedGalleries/myGallery1/images/myImage1",
			Expected: &SharedGalleryImageId{
				GalleryName: "myGallery1",
				ImageName:   "myImage1",
			},
		},

		{
			// upper-cased
			Input: "/SHAREDGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SharedGalleryImageID(v.Input)
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
