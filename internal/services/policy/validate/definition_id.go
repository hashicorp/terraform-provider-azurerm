// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
)

func PolicyDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PolicyDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Policy Definition ID: %+v", k, err))
		return
	}

	return warnings, errors
}
