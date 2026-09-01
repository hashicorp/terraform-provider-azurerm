// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 63),
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9]([a-z0-9-]+[a-z0-9])?$`), "must only contains numbers, lowercase characters and '-'"),
	)(i, k)
}
