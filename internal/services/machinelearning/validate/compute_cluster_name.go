// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// The portal says: It can include letters, digits and dashes. It must start with a letter, end with a letter or digit, and be between 3 and 32 characters in length.
// If you provide invalid name, the rest api will return an error with the following regex.
func ComputeClusterName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{1,30}[a-zA-Z0-9]$`), "must be between 3 and 32 characters, may only include alphanumeric characters and '-' and must start with a letter, end with a letter or digit")(i, k)
}
