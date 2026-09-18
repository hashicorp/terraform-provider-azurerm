// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func OrganizationOrganizationID(i interface{}, k string) ([]string, []error) {
	// Ensure the string follows the desired format.
	// Regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$
	// The negative lookahead (?!-) is not supported in Go's standard regexp package, so it is
	// expressed equivalently by requiring the first character to be [a-z0-9].
	return validation.StringMatch(regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,47}[a-z0-9]$`), "iot central organizationId is invalid, regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$")(i, k)
}
