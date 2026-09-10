// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagementLockName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`[A-Za-z0-9-_]`), "can only consist of alphanumeric characters, dashes and underscores"),
		validation.StringLenBetween(0, 259),
	)(v, k)
}
