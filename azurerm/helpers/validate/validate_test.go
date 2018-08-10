package validate

import "testing"

func TestUrlWithScheme(t *testing.T) {
	validSchemes := []string{"example"}
	testCases := []struct {
		Url             string
		ShouldHaveError bool
	}{
		{
			Url:             "example://mysite.com",
			ShouldHaveError: false,
		},
		{
			Url:             "http://mysite.com",
			ShouldHaveError: true,
		},
		{
			Url:             "example://",
			ShouldHaveError: true,
		},
		{
			Url:             "example://validhost",
			ShouldHaveError: false,
		},
	}

	t.Run("TestUrlWithScheme", func(t *testing.T) {
		for _, v := range testCases {
			_, errors := UrlWithScheme(validSchemes)(v.Url, "field_name")

			hasErrors := len(errors) > 0
			if v.ShouldHaveError && !hasErrors {
				t.Fatalf("Expected an error but didn't get one for %q", v.Url)
				return
			}

			if !v.ShouldHaveError && hasErrors {
				t.Fatalf("Expected %q to return no errors, but got some %+v", v.Url, errors)
				return
			}
		}
	})
}
