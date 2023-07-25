// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

// NormalizeLocation will be deprecated in the near future, use `location.Normalize()` instead.
func NormalizeLocation(input interface{}) string {
	loc := input.(string)
	return location.Normalize(loc)
}

// SchemaResourceGroupNameDiffSuppress will be deprecated in the near future
// use `commonschema.ResourceGroupName()` instead
func SchemaResourceGroupNameDiffSuppress() *pluginsdk.Schema {
	// @tombuildsstuff: this function should no longer be used, existing resources will need to be worked
	// through (and this switched out for `commonschema.ResourceGroupName()`) once verifying these now pull
	// this value from the Resource ID rather than the API Response.

	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppress.CaseDifference,
		ValidateFunc:     resourcegroups.ValidateName,
	}
}

func ValidateResourceID(i interface{}, k string) (warnings []string, errors []error) {
	// ValidateResourceID should only be used when a more specific Resource ID validation function
	// is unavailable.
	// If in doubt, prefer a validation function that supports multiple types of validation functions
	// rather than using this function
	// e.g. `validation.Any(commonids.ValidateSubnetID, commonids.ValidateVirtualNetworkID)`
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAzureResourceID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
	}

	return warnings, errors
}

// Deprecated: use a more specific Resource ID validator instead, however note that empty strings should not
// be allowed as a validation value. Rather than specifying an empty string, users can omit the field to use
// an unset value.
func ValidateResourceIDOrEmpty(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		return
	}

	return ValidateResourceID(i, k)
}
