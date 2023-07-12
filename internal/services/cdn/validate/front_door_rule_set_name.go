// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorRuleSetName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z][\da-zA-Z]{0,59}$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 1 and 60 characters in length, begin with a letter and may contain only letters and numbers, got %q`, k, v)}
	}

	return nil, nil
}
