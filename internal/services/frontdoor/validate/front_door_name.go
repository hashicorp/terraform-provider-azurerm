// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FrontDoorName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`(^[\da-zA-Z])([-\da-zA-Z]{3,61})([\da-zA-Z]$)`), "must be between 5 and 63 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens")(i, k)
}
