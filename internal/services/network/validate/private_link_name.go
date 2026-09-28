// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PrivateLinkName(i interface{}, k string) ([]string, []error) {
	// The name attribute rules per the Nat Gateway service team are (Friday, October 18, 2019 4:20 PM):
	// 1. Must not be empty.
	// 2. Must be between 1 and 80 characters.
	// 3. The attribute must:
	//    a) begin with a letter or number
	//    b) end with a letter, number or underscore
	//    c) may contain only letters, numbers, underscores, periods, or hyphens.
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z\d]$|^[a-zA-Z\d][a-zA-Z\d-_.]{0,78}[a-zA-Z\d_]$`),
		"must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, periods, hyphens or underscores",
	)(i, k)
}
