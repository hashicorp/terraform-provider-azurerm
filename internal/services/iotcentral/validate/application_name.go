// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ApplicationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[a-z\d][a-z\d-]{0,61}[a-z\d]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain lowercase alphanumeric characters and dashes, length between 2-63, and must start and end with an alphanumeric character", k))
	}
	return warnings, errors
}
