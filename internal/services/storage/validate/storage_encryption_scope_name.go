// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageEncryptionScopeName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[0-9a-zA-Z]{4,63}$"), "must be alphanumeric, and between 4 to 63 characters")(v, k)
}
