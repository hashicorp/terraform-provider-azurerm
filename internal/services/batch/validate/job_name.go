// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// JobName validates the name of a Batch job
func JobName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\w-]{1,64}$`), "can only contain any combination of alphanumeric characters along with dash (-) and underscore (_). The name must be from 1 through 64 characters long")(v, k)
}
