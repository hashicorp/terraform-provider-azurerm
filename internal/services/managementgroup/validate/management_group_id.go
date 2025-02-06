// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
)

func ManagementGroupID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ManagementGroupID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a management group id: %v", k, err))
		return
	}

	return
}

func TenantScopedManagementGroupID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.TenantScopedManagementGroupID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a management group id: %v", k, err))
		return
	}

	return
}
