// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IotSecuritySolutionName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([-a-zA-Z0-9_.])+$`), "can only contain letter, digit, '-', '.' or '_'")(input, key)
}
