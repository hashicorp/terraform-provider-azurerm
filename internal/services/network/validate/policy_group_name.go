// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PolicyGroupName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z][a-zA-Z.\-_]{0,78}[a-zA-Z_]{0,1}$`),
		"must be between 1 and 80 characters in length, must begin with a letter, end with a letter or underscore and contain only letters, periods, underscores and hyphens",
	)(i, k)
}
