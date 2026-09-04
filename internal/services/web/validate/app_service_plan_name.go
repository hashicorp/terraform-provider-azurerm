// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AppServicePlanName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z-_]{1,60}$`), "may only contain alphanumeric characters, dashes and underscores up to 60 characters in length")(v, k)
}
