// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func KeyVaultChildID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.ParseNestedItemID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Key Vault Child resource id: %v", k, err))
	}

	return warnings, errors
}

func KeyVaultChildIDWithOptionalVersion(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.ParseOptionallyVersionedNestedItemID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Key Vault Child resource id: %v", k, err))
	}

	return warnings, errors
}
