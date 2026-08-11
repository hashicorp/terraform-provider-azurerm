// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 63),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z-_]`), "must begin with a letter, `-` or `_`"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-_]+$`), "must only contain numbers, characters and `-`, `_`"),
	)(i, k)
}
