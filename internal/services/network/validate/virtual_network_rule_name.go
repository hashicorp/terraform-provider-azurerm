// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VirtualNetworkRuleName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 128),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-]*$`), "can only contain alphanumeric characters and hyphens"),
		validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
		validation.StringDoesNotMatch(regexp.MustCompile(`^[0-9-]`), "cannot start with a number or hyphen"),
	)(v, k)
}
