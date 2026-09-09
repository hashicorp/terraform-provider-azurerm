// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SqlAdministratorLoginName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d-]{0,127}$`), "can contain only letters, numbers, or dashes, must start with a letter, The value must be between 1 and 128 characters long")(i, k)
}
