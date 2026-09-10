// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LinuxComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name cannot exceed 64 characters in length
	return LinuxComputerName(i, k, 64, false)
}

func LinuxComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name prefix cannot exceed 58 characters in length
	return LinuxComputerName(i, k, 58, true)
}

func LinuxComputerName(i interface{}, k string, maxLength int, allowDashSuffix bool) (warnings []string, errors []error) {
	validator := validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringLenBetween(1, maxLength),
		validation.StringDoesNotMatch(regexp.MustCompile(`^_`), "cannot begin with an underscore"),
		validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "cannot end with a period"),
		// Linux host name cannot contain the following characters
		validation.StringDoesNotContainAny(`\/"[]:|<>+=;,?*@&~!#$%^()_{}'`),
	)
	if !allowDashSuffix {
		validator = validation.All(validator, validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a dash"))
	}

	return validator(i, k)
}
