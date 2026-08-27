// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ExportName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespace"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "may only contain letters and digits"),
		validation.StringLenBetween(3, 24),
	)(v, k)
}
