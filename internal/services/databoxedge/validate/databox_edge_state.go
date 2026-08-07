// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataboxEdgeState(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\s\S]{1,30}$`), "must be between 1 and 30 characters in length")(v, k)
}
