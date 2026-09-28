// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DashboardName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(0, 160),
		validation.StringMatch(regexp.MustCompile(`^[-\w]+$`), "may only contain alphanumeric and hyphen characters"),
	)(v, k)
}
