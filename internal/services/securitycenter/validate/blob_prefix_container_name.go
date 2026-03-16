// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func BlobPrefixContainerName(v interface{}, k string) (warnings []string, errors []error) {
	value := strings.Split(v.(string), "/")[0]

	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only lowercase alphanumeric characters and hyphens allowed in %q container name: %q", k, value))
	}
	if len(value) > 63 {
		errors = append(errors, fmt.Errorf("%q container name must be less than 63 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q container name cannot begin with a hyphen: %q", k, value))
	}

	return warnings, errors
}
