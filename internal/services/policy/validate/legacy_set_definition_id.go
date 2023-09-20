// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
)

func PolicySetDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PolicySetDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Policy Set Definition ID: %+v", v, err))
		return
	}

	return warnings, errors
}
