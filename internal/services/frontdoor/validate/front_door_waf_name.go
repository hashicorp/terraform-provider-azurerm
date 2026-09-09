// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FrontDoorWAFName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`(^[a-zA-Z])([\da-zA-Z]{0,127})$`), "must be between 1 and 128 characters in length, must begin with a letter and may only contain letters and numbers")(i, k)
}
