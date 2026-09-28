// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RuleName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_0-9.-]+$`), "only word characters, numbers, underscores, periods, and hyphens allowed"),
		validation.StringLenBetween(1, 80),
		validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9_]$`), "must end with a word character, number, or underscore"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]`), "must start with a word character or number"),
	)(v, k)
}
