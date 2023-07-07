// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorOriginGroupName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[\da-zA-Z][-\da-zA-Z]{0,88}[\da-zA-Z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 2 and 90 characters in length, begin with a letter or number, end with a letter or number and may contain only letters, numbers and hyphens, got %q`, k, i))
	}

	return nil, nil
}
