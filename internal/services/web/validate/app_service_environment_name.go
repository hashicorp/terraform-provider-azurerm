// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AppServiceEnvironmentName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{0,61}[0-9a-zA-Z]$`), "may only contain alphanumeric characters and dashes up to 63 characters in length, and must start and end in an alphanumeric")(v, k)
}
