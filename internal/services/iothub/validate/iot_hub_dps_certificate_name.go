// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IoTHubDpsCertificateName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringIsNotEmpty,
		validation.StringLenBetween(1, 64),
		validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z-._]+$`), "may only contain alphanumeric characters or the following: -._"),
	)(v, k)
}
