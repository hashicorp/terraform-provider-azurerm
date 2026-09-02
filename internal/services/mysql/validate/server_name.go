// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ServerName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z][-0-9a-z]{1,61}[0-9a-z]$`), "can contain only lowercase letters, numbers, and '-', but can't start or end with '-', and must be at least 3 characters and no more than 63 characters long")(i, k)
}
