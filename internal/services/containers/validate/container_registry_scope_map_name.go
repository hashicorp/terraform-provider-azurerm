// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ContainerRegistryScopeMapName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\-]+$`), "alpha numeric characters and hyphens only are allowed"),
		validation.StringLenBetween(5, 49),
	)(v, k)
}
