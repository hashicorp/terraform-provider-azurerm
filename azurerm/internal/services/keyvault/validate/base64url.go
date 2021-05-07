package validate

import (
	"encoding/base64"
	"fmt"
)

// StringIsRawBase64Url is a ValidateFunc that ensures a string can be parsed as RawBase64Url
func StringIsRawBase64Url(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := base64.RawURLEncoding.DecodeString(v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a Raw base64 URL string, got %v", k, v))
	}

	return warnings, errors
}
