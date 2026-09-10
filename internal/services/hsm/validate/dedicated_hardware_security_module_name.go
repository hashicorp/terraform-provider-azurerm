// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DedicatedHardwareSecurityModuleName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$`), "must be between 3 and 24 alphanumeric characters. It must begin with a letter, end with a letter or digit"),
		// No consecutive hyphens
		validation.StringDoesNotMatch(regexp.MustCompile("(--)"), "must not contain any consecutive hyphens"),
	)(i, k)
}
