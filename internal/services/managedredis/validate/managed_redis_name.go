// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagedRedisClusterName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 63),
		validation.StringDoesNotMatch(regexp.MustCompile(`--`), "must not contain any consecutive hyphens"),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`), "can only contain letters, numbers and hyphens. The first and last characters must each be a letter or a number"),
	)(i, k)
}
