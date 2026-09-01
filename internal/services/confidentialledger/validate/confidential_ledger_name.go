// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ConfidentialLedgerName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 32),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-`), "may not start with a dash"),
		validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "may not end with a dash"),
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		validation.StringMatch(regexp.MustCompile(`^[^\-][A-Za-z0-9\-]{1,33}[^\-]$`), "may only contain alphanumeric characters and dash"),
	)(v, k)
}
