// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	validatehelper "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorSecurityPolicyName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validatehelper.RegExHelper(i, k, `^[\da-zA-Z](?:[-\da-zA-Z]*[\da-zA-Z])?$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q must begin and end with an alphanumeric character, and may contain only alphanumeric characters and hyphens", k))
	}

	return nil, nil
}
