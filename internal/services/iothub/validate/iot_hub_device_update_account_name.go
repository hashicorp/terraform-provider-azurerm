// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IotHubDeviceUpdateAccountName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 24),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`), "must start with an alphanumeric, may only contain alphanumeric characters and dashes, and consecutive dashes (-) are not allowed"),
	)(v, k)
}
