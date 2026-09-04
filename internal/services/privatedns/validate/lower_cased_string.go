// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// LowerCasedString validates that the string is lower-cased
func LowerCasedString(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringDoesNotMatch(regexp.MustCompile(`[\p{Lu}\p{Lt}]`), "must be a lower-cased string"),
		validation.StringDoesNotContainAny(" "),
	)(i, k)
}
