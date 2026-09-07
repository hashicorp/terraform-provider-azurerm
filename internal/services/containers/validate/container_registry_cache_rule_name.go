// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ContainerRegistryCacheRuleName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`), "alpha numeric characters optionally separated by '-' only are allowed"),
		validation.StringLenBetween(5, 49),
	)(v, k)
}
