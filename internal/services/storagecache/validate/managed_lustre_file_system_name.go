// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagedLustreFileSystemName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`), "can contain alphanumeric characters, hyphens and underscores and start and end with alphanumeric and has to be between 2 and 80 characters")(i, k)
}
