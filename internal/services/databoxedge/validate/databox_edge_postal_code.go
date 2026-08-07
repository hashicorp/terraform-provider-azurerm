// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataboxEdgePostalCode(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([\d]{5})((-)([\d]{4}))?$`), "must consist of 5 digits unless it is in the ZIP+4 format then it must consist of five digits, a hyphen, then four digits")(v, k)
}
