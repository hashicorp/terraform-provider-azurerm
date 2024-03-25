// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func UserType(i interface{}, k string) (warnings []string, errors []error) {
	userType, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	err := validateUserType(userType)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func validateUserType(userType string) error {
	validUserTypes := []string{"Group", "ServicePrincipal", "Email"}

	for _, t := range validUserTypes {
		if userType == t {
			return nil
		}
	}

	return fmt.Errorf("iot central userType %q is invalid, expected one of: %s", userType, strings.Join(validUserTypes, ", "))
}
