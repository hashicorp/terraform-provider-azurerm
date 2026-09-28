// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LogAnalyticsWorkspaceName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$"), "can only contain alphabet, number, and '-' character, and cannot use '-' as the start and end of the name"),
		validation.StringLenBetween(4, 63),
	)(v, k)
}
