// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
)

func EmailUserType(input interface{}, key string) (warnings []string, errors []error) {
	err := validateUserType(input, key, "Email")
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func ServicePrincipalUserType(input interface{}, key string) (warnings []string, errors []error) {
	err := validateUserType(input, key, "ServicePrincipal")
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func GroupUserType(input interface{}, key string) (warnings []string, errors []error) {
	err := validateUserType(input, key, "Group")
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func validateUserType(input interface{}, key string, expectedType string) error {
	value, ok := input.(string)
	if !ok {
		return fmt.Errorf("expected %s to be a string", key)
	}

	if value == expectedType {
		return nil
	}

	return fmt.Errorf("expected %s to be %s but got %s", key, expectedType, value)
}
