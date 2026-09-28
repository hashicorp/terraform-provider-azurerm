// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FirewallRuleName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,}[a-zA-Z0-9_]$`), "must consist of letters, digits, underscores, periods and hyphens. The first character must be a letter or digit, and the last character must be a letter, a digit or an underscore")(v, k)
}
