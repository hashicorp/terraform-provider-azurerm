// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ResourceAnchorName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{N}-]*$`), "must contain only letters , numbers and hyphens"),
		validation.StringLenBetween(0, 24),
	)(i, k)
}
