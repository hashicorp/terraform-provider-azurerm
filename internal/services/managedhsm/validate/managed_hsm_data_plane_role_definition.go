// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ManagedHSMDataPlaneRoleDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	if _, err := parse.ManagedHSMDataPlaneRoleDefinitionID(v, nil); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v, err))
		return warnings, errors
	}

	return warnings, errors
}
