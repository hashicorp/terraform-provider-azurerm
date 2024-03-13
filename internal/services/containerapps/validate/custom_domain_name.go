// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
)

// ContainerAppCustomDomainId checks that 'input' can be parsed as a Container App ID
func ContainerAppCustomDomainId(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.ContainerAppCustomDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
