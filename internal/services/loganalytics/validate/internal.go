// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsGenericName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(4, 63),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`), "expected value does not match regular expression"),
	)(i, k)
}
