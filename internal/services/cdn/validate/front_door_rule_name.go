// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FrontDoorRuleName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][\da-zA-Z]{0,259}$`), "must be between 1 and 260 characters in length, begin with a letter and may contain only letters and numbers")(v, k)
}
