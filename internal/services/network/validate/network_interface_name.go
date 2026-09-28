// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func NetworkInterfaceName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][0-9a-zA-Z_.-]{0,62}[0-9a-zA-Z_]$`), "must be between 2 and 64 characters in length, must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods or hyphens")(v, k)
}
