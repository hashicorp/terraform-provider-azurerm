// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AlertProcessingRuleName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([a-zA-Z\d])[a-zA-Z\d-_]*$`), "should begin with a letter or number, contain only letters, numbers, underscores and hyphens")(i, k)
}
