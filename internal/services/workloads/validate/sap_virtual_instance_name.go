// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SAPVirtualInstanceName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Z][A-Z0-9][A-Z0-9]$`), "must have three characters. The first character must be alphabetic. The second and third characters must be alphanumerical. All alphabetic characters should be uppercase")(v, k)
}
