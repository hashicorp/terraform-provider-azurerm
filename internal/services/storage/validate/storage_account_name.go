// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageAccountName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`\A([a-z0-9]{3,24})\z`), "can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long")(v, k)
}
