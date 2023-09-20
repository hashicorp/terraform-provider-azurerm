// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
)

func PolicyAssignmentID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PolicyAssignmentID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Policy Assignment ID: %+v", k, err))
		return
	}

	return warnings, errors
}
