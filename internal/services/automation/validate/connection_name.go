// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ConnectionName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\w\-]{1,128}$`), "contain only letters, numbers hyphens and underscore. The value must be between 1 and 128 characters long")(i, k)
}
