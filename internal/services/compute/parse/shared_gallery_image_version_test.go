// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SharedGalleryImageVersionId{}

func TestSharedGalleryImageVersionIDFormatter(t *testing.T) {
	actual := NewSharedGalleryImageVersionID("myGallery1", "myImage1", "latest").ID()
	expected := "/sharedGalleries/myGallery1/images/myImage1/versions/latest"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSharedGalleryImageVersionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SharedGalleryImageVersionId
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
			// missing versions
			Input: "/sharedGalleries/myGallery1/images/myImage1",
			Error: true,
		},

		{
			// missing value for versions
			Input: "/sharedGalleries/myGallery1/images/versions",
			Error: true,
		},

		{
			// valid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/latest",
			Expected: &SharedGalleryImageVersionId{
				GalleryName: "myGallery1",
				ImageName:   "myImage1",
				Version:     "latest",
			},
		},

		{
			// valid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/1.2.3",
			Expected: &SharedGalleryImageVersionId{
				GalleryName: "myGallery1",
				ImageName:   "myImage1",
				Version:     "1.2.3",
			},
		},

		{
			// invalid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/notTheLatest",
			Error: true,
		},

		{
			// invalid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/1.2",
			Error: true,
		},

		{
			// invalid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/1.2.",
			Error: true,
		},

		{
			// invalid
			Input: "/sharedGalleries/myGallery1/images/myImage1/versions/1.two.3",
			Error: true,
		},

		{
			// upper-cased
			Input: "/SHAREDGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1/VERSIONS/1.2.3",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SharedGalleryImageVersionID(v.Input)
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

		if actual.Version != v.Expected.Version {
			t.Fatalf("Expected %q but got %q for Versions", v.Expected.Version, actual.Version)
		}
	}
}
