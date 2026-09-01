// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// IsUUIDOrEmpty is a ValidateFunc that ensures a string can be parsed as UUID or is empty
func IsUUIDOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	return validation.Any(
		validation.StringIsEmpty,
		validation.IsUUID,
	)(i, k)
}
