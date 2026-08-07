// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataboxEdgePhoneNumber(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([\+])?(\d{1,2})?[\s-]?\(?(\d{3})\)?[\s-]?(\d{3})[-](\d{4})$`), "may contain parentheses, hyphens, plus sign or digits only, must be in a valid 10 digit phone number format with country code being optional(e.g. 123 555-6789)")(v, k)
}
