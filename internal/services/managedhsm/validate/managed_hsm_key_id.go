// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func OptionallyVersionedNestedItemID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if _, err := parse.ParseOptionallyVersionedNestedItemID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %s", v, err))
		return warnings, errors
	}

	return warnings, errors
}
