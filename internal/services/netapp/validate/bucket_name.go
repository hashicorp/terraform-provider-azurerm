// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// BucketName validates that the given value is a valid Azure NetApp Files Bucket name.
// The bucket name is S3-compatible: 3-63 characters, DNS-compliant, lowercase
// letters/numbers/hyphens/periods, must start and end with a letter or number,
// must not contain consecutive periods or "."- / "-." sequences, and must not
// look like an IPv4 address.
func BucketName(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 63 characters in length", k))
		return warnings, errors
	}

	// Allowed characters only.
	if !regexp.MustCompile(`^[a-z0-9.\-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only lowercase letters, numbers, hyphens or periods", k))
		return warnings, errors
	}

	// Must start with a lowercase letter or a digit.
	if !regexp.MustCompile(`^[a-z0-9]`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must start with a lowercase letter or number", k))
		return warnings, errors
	}

	// Must end with a lowercase letter or a digit.
	if !regexp.MustCompile(`[a-z0-9]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must end with a lowercase letter or number", k))
		return warnings, errors
	}

	// No consecutive periods, no ".-" or "-." sequences.
	if regexp.MustCompile(`\.\.|\.\-|\-\.`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must not contain consecutive periods or %q / %q sequences", k, ".-", "-."))
		return warnings, errors
	}

	// Must not look like an IPv4 address.
	if regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must not be formatted as an IPv4 address", k))
		return warnings, errors
	}

	return warnings, errors
}
