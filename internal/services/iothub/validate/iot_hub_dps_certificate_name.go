// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func IoTHubDpsCertificateName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q may contain at most 64 characters", k))
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-._]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters or the following: -._", k))
	}

	return warnings, errors
}
