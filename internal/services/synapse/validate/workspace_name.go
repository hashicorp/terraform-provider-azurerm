// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func WorkspaceName(i interface{}, k string) ([]string, []error) {
	// The name attribute rules are :
	// 1. can contain only lowercase letters, numbers or hyphens
	// 2. must start and end with a lowercase letter or number
	// 3. must not contain the string '-ondemand'
	// 4. The value must be between 1 and 50 characters long
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,48}[a-z0-9])?$`), "must start and end with a letter or number, can contain only lowercase letters, numbers or hyphens, and be between 1 and 50 characters long"),
		validation.StringDoesNotMatch(regexp.MustCompile(`-ondemand`), "must not contain the string '-ondemand'"),
	)(i, k)
}
