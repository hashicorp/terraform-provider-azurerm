// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SnapshotName(v interface{}, k string) (warnings []string, errors []error) {
	// a-z, A-Z, 0-9, _ and -. The max name length is 80
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9_-]+$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s can only contain alphanumeric characters and underscores", k))
	}

	length := len(value)
	if length > 80 {
		errors = append(errors, fmt.Errorf("%s can be up to 80 characters, currently %d", k, length))
	}

	return warnings, errors
}
