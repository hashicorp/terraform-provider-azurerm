// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageTargetName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[-0-9a-zA-Z_]{1,31}$`), "can contain alphanumeric characters, dashes and underscores and has to be between 1 and 31 characters")(i, k)
}
