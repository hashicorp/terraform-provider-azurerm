// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorCustomDomainName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z0-9][a-zA-Z0-9-]{0,258}[a-zA-Z0-9]$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 2 and 260 characters in length, must begin with a letter or number, end with a letter or number and contain only letters, numbers and hyphens, got %q`, k, v)}
	}

	return nil, nil
}
