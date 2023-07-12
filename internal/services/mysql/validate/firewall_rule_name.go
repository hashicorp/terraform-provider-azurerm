// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FirewallRuleName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[a-zA-Z0-9-_]{1,128}$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only letters, numbers, underscores and dashes. It must be less than or equal to 128 characters", k))
	}

	return nil, nil
}
