// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func HubRouteTableName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[^<>%&:?/+]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must not contain characters from %q", k, "<>&:?/+%"))
	}

	return warnings, errors
}
