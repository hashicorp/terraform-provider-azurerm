// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationGroupName(i interface{}, k string) (warnings []string, errors []error) {
	// Application Group name can be 3-64 characters in length
	return validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringLenBetween(3, 64),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9._-]+$`), "may only contain alphanumeric characters, dots, dashes and underscores"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]`), "must begin with an alphanumeric character"),
		validation.StringMatch(regexp.MustCompile(`\w$`), "must end with an alphanumeric character or underscore"),
	)(i, k)
}
