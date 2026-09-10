// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagedRedisDatabaseGeoreplicationGroupName(val interface{}, argName string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 63),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "can only contain letters, numbers, and hyphens"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]`), "must start with a letter or number"),
		validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9]$`), "must end with a letter or number"),
		validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain consecutive hyphens"),
	)(val, argName)
}
