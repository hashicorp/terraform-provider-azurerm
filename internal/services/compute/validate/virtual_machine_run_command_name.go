// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VirtualMachineRunCommandName(i interface{}, k string) (warnings []string, errors []error) {
	// Run Command name can be 1-80 characters in length
	return validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringLenBetween(1, 80),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9._-]+$`), "may only contain alphanumeric characters, dots, dashes and underscores"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]`), "must begin with an alphanumeric character"),
		validation.StringMatch(regexp.MustCompile(`\w$`), "must end with an alphanumeric character or underscore"),
	)(i, k)
}
