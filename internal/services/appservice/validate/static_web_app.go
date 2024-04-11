// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func StaticWebAppName(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,60}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and up to 60 characters in length", k))
	}

	return warnings, errors
}

func StaticWebAppPassword(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return
	}

	if len(value) < 8 {
		errors = append(errors, fmt.Errorf("the password should be at least eight characters long."))
	}

	if matched := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("the password must contain at least one symbol"))
	}

	if matched := regexp.MustCompile(`[a-z]`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("the password must contain at least one lower case letter"))
	}

	if matched := regexp.MustCompile(`[A-Z]`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("the password must contain at least one upper case letter"))
	}

	if matched := regexp.MustCompile(`[0-9]`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("the password must contain at least one number"))
	}

	if matched := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("the password must contain at least one symbol"))
	}

	return
}
