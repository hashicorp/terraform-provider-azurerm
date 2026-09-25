// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FileSystemName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,13}[a-zA-Z0-9]$`), "must be between 2 and 15 characters in length, must not begin or end with a hyphen and may only contain alphanumeric characters and hyphens")(v, k)
}
