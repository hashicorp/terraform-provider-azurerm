// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func PolicyGroupName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z][a-zA-Z.-_]{0,78}[a-zA-Z_]{0,1}$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 1 and 80 characters in length, must begin with a letter, end with a letter or underscore and contain only letters, periods, underscores and hyphens, got %q`, k, v)}
	}

	return nil, nil
}
