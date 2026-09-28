// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// The portal says: The workspace name must be between 3 and 33 characters. The name may only include alphanumeric characters and '-'.
// If you provide invalid name, the rest api will return an error with the following regex.
func WorkspaceName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][\w-]{2,32}$`), "must be between 3 and 33 characters, and may only include alphanumeric characters and '-'")(i, k)
}
