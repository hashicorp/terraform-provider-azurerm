// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func WebAppName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,60}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and up to 60 characters in length", k))
	}

	return warnings, errors
}

func ContainerizedFunctionAppName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-z])[a-z0-9-]{0,58}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 60 characters", k))
		return
	}

	return
}
