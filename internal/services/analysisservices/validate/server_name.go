// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ServerName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z][0-9a-z]{2,62}$`), "must begin with a letter, be lowercase alphanumeric, and be between 3 and 63 characters in length")(v, k)
}
