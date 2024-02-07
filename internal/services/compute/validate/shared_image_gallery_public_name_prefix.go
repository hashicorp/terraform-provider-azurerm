// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SharedImageGalleryPrefix(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9]{5,16}$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be 5 to 16 characters long, and can only contain alphanumeric", k))
	}

	return warnings, errors
}
