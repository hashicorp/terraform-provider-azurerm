// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DatadogMonitorsName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_-]{2,32}$`), "must be between 2 and 32 characters in length, can only contain alphanumeric characters, underscore and hyphen symbols")(i, k)
}
