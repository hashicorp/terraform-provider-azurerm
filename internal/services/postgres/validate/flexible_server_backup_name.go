// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerBackupName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if !regexp.MustCompile(`^[-\w\._]+$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid backup name, got %v", k, v))
		return warnings, errors
	}
	return warnings, errors
}
