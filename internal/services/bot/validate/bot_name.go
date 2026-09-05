// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BotName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(4, 42),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`), "must start with a letter or digit and may only contain alphanumeric characters, underscores and dashes"),
	)(i, k)
}
