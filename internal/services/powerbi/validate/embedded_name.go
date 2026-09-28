// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func EmbeddedName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]{3,63}$`), "must be between 4 and 64 characters in length and starts with a letter and contains only lowercase letters or numbers")(v, k)
}
