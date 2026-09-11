// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AccountName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[-_\da-zA-Z]{3,64}$`), "must be between 3 and 64 characters in length and contains only letters, numbers, underscore or hyphens")(v, k)
}
