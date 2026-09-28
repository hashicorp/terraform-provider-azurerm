// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IotHubSharedAccessPolicyName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9!._-]{1,64}`), "must not be empty, and must not exceed 64 characters in length, and can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens")(i, k)
}
