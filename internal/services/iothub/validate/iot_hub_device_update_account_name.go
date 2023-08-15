// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func IotHubDeviceUpdateAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 3 || len(value) > 24 {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 24 characters long", k))
	}

	if matched := regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must start with an alphanumeric, may only contain alphanumeric characters and dashes, and consecutive dashes (-) are not allowed", k))
	}

	return warnings, errors
}
