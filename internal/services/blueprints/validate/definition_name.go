// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DefinitionName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-_]{1,48}$`), "can include letters, numbers, underscores or dashes. Spaces and other special characters are not allowed")(i, k)
}
