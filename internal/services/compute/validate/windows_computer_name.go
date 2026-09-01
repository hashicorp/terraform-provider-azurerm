// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func WindowsComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name cannot be more than 15 characters long
	return windowsComputerName(i, k, 15, false)
}

func WindowsComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name prefix cannot be more than 9 characters long
	return windowsComputerName(i, k, 9, true)
}

func windowsComputerName(i interface{}, k string, maxLength int, allowDashSuffix bool) (warnings []string, errors []error) {
	validator := validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringLenBetween(1, maxLength),
		// A windows computer name can only contain alphanumeric characters and hyphens
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "may only contain alphanumeric characters and dashes"),
		// Windows computer name cannot contain only numbers
		validation.StringDoesNotMatch(regexp.MustCompile(`^\d+$`), "cannot contain only numbers"),
	)
	if !allowDashSuffix {
		validator = validation.All(validator, validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a dash"))
	}

	return validator(i, k)
}
