package validate

import (
	"fmt"
	"strings"
)

func SharedGalleryImageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if strings.Contains(strings.ToLower(v), "/sharedgalleries/") {
		errors = append(errors, fmt.Errorf("%q is missing the %q segment, got %q", key, "/sharedGalleries/", v))
	}

	return
}
