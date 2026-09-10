// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FirewallRuleName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^<>*%&:\\/?]{0,127}[^.<>*%&:\\/?]$`), "can't contain '<,>,*,%,&,:,\\,/,?', can't end with '.', and must be between 1 and 128 characters long")(i, k)
}
