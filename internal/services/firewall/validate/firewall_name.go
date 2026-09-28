// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FirewallName(v interface{}, k string) ([]string, []error) {
	// From the Portal:
	// The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z._-]{0,}[0-9a-zA-Z_])?$`), "must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens")(v, k)
}
