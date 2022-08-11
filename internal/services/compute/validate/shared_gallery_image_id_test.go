package validate

import "testing"

func TestSharedGalleryImageID(t *testing.T) {
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
			// missing SharedGalleries
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SharedGalleries
			Input: "/SharedGalleries/",
			Valid: false,
		},

		{
			// missing images
			Input: "/SharedGalleries/myGallery1/",
			Valid: false,
		},

		{
			// missing value for images
			Input: "/SharedGalleries/myGallery1/images/",
			Valid: false,
		},

		{
			// valid
			Input: "/SharedGalleries/myGallery1/images/myImage1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SHAREDGALLERIES/MYGALLERY1/IMAGES/MYIMAGE1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SharedGalleryImageID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t: %+v", tc.Valid, valid, errors)
		}
	}
}
