// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerFirewallRuleName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 128),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-_]+$`), "must only contains numbers, characters and `-`, `_`"),
	)(i, k)
}
