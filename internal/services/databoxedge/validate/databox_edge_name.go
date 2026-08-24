// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataboxEdgeName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\da-zA-Z][-\da-zA-Z]{1,22}[\da-zA-Z]$`), "must be between 3 and 24 characters in length, begin and end with an alphanumeric character, can only contain alphanumeric characters and hyphens")(v, k)
}
