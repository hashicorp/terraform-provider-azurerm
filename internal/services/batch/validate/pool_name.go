// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// PoolName validates the name of a Batch pool
func PoolName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_-]+$`), "any combination of alphanumeric characters including hyphens and underscores are allowed"),
		validation.StringLenBetween(1, 64),
	)(v, k)
}
