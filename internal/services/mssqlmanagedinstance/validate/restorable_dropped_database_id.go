// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func RestorableDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.RestorableDroppedDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a MsSql Managed Instance Restorable Database resource id: %v", k, err))
	}

	return warnings, errors
}
