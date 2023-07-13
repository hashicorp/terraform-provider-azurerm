// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func DatabricksVirtualNetworkPeeringName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	//  and must be between 1 and 80 characters in length
	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z\d][a-zA-Z\d._-]{0,78}[a-zA-Z\d_]$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 2 and 80 characters in length, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens, got %q`, k, v)}
	}

	return nil, nil
}
