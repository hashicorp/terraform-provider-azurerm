// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PublicIpDomainNameLabel(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9-]{1,61}[a-z0-9]$`), "must contain only lowercase alphanumeric characters, numbers and hyphens. It must start with a letter, end only with a number or letter and not exceed 63 characters in length")(v, k)
}
