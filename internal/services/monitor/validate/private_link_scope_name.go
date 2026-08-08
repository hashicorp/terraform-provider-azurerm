// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PrivateLinkScopeName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9()-_.]*[a-zA-Z0-9_-]$`), "only allows alphanumeric characters, periods, underscores, hyphens and parenthesis and cannot end in a period"),
	)(i, k)
}
