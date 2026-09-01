// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CustomIpPrefixName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_])$`), "must be between 1 and 80 characters in length, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods or hyphens")(i, k)
}
