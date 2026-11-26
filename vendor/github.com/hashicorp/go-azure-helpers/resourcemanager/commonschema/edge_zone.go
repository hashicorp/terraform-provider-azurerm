// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// EdgeZoneComputed returns the schema for an Edge Zone which is Computed
func EdgeZoneComputed() *schema.Schema {
	return &schema.Schema{
		Description: "The Edge Zone within the Azure Region where this resource is deployed.",
		Type:        schema.TypeString,
		Computed:    true,
	}
}

// EdgeZoneOptional returns the schema for an Edge Zone which is Optional
func EdgeZoneOptional() *schema.Schema {
	return &schema.Schema{
		Description:      "The Edge Zone within the Azure Region to deploy to.",
		Type:             schema.TypeString,
		Optional:         true,
		ValidateFunc:     validation.StringIsNotEmpty,
		StateFunc:        edgezones.StateFunc,
		DiffSuppressFunc: edgezones.DiffSuppressFunc,
	}
}

// EdgeZoneOptionalForceNew returns the schema for an Edge Zone which is both Optional and ForceNew
func EdgeZoneOptionalForceNew() *schema.Schema {
	return &schema.Schema{
		Description:      "The Edge Zone within the Azure Region to deploy to.",
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateFunc:     validation.StringIsNotEmpty,
		StateFunc:        edgezones.StateFunc,
		DiffSuppressFunc: edgezones.DiffSuppressFunc,
	}
}
