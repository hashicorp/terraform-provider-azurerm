// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VirtualNetworkRuleName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(2, 64),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-\._]*$`), "can only contain alphanumeric characters, underscores, periods and hyphens"),
		validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
		validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "cannot end with a period"),
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\._-]`), "cannot start with a period, underscore or hyphen"),
	)(v, k)
}
