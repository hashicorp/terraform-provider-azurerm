// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageTableName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringNotInSlice([]string{"table"}, false),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{2,62}$`), "cannot begin with a numeric character, only alphanumeric characters are allowed and must be between 3 and 63 characters long"),
	)(v, k)
}
