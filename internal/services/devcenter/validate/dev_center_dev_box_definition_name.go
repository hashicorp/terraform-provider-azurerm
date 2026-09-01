// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DevCenterDevBoxDefinitionName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-_.]{2,62}$"), "must start with an alphanumeric character, may contain alphanumeric characters, dashes, underscores or periods and must be between 3 and 63 characters long")(i, k)
}
