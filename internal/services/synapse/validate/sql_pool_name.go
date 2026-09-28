// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SqlPoolName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^<>*%&:\\\/?@-]{0,59}[^\s.<>*%&:\\\/?@-]$`), "can contain only letters, numbers or underscore, The value must be between 1 and 60 characters long")(i, k)
}
