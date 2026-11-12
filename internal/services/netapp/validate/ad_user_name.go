// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// ActiveDirectoryUserName validates that the given value is a valid Azure NetApp Account AD user name.
// User names added to a NetApp managed Active Directory must not exceed 255 characters in length,
// and not contain the characters '@' or '\' (backslash).
func ActiveDirectoryUserName(v any, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringDoesNotContainAny(`@\`),
	)(v, k)
}
