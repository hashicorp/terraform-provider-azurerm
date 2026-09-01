// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagedHardwareSecurityModuleName(i interface{}, k string) ([]string, []error) {
	// The name attribute rules are :
	// Must be a 3-24 character string, containing only 0-9, a-z. A-Z, and -
	// The name must begin with a letter, end with a letter or digit, and not contain consecutive hyphens.
	msg := "must begin with a letter, end with a letter or number, contain only alphanumeric characters. The value must be between 3 and 24 characters long"
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`--`), msg),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][-a-zA-Z\d]{1,22}[a-zA-Z\d]$`), msg),
	)(i, k)
}
