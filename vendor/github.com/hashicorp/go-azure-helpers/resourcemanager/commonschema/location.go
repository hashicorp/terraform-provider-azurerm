// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Location() *schema.Schema {
	return LocationWithDescription("The Azure Region in which to provision the resource.")
}

func LocationWithDescription(description string) *schema.Schema {
	return &schema.Schema{
		Description:      description,
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateFunc:     location.EnhancedValidate,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}

func LocationOptional() *schema.Schema {
	return LocationOptionalWithDescription("The Azure Region in which to provision the resource.")
}

func LocationOptionalWithDescription(description string) *schema.Schema {
	return &schema.Schema{
		Description:      description,
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}

func LocationComputed() *schema.Schema {
	return LocationComputedWithDescription("The Azure Region in which the resource exists.")
}

func LocationComputedWithDescription(description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeString,
		Computed:    true,
	}
}

func LocationWithoutForceNew() *schema.Schema {
	return LocationWithoutForceNewWithDescription("The Azure Region in which to provision the resource.")
}

func LocationWithoutForceNewWithDescription(description string) *schema.Schema {
	return &schema.Schema{
		Description:      description,
		Type:             schema.TypeString,
		Required:         true,
		ValidateFunc:     location.EnhancedValidate,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}
