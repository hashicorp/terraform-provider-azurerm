// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// BucketPath validates that the given value is an absolute POSIX-style path used
// as the mount path inside a NetApp Files bucket.
func BucketPath(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if !regexp.MustCompile(`^/`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be an absolute POSIX-style path starting with %q", k, "/"))
	}

	if regexp.MustCompile(`\\`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must not contain backslashes", k))
	}

	return warnings, errors
}
