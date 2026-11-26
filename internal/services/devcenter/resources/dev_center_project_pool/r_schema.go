// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Project Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterProjectPoolName,
		},

		"location": commonschema.Location(),

		"dev_center_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&pools.ProjectId{}, "Dev Center Project"),

		"dev_box_definition_name": {
			Description:  "The name of the Dev Center Dev Box Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"local_administrator_enabled": {
			Description: "Whether owners of Dev Boxes in the Dev Center Project Pool are added as local administrators on the Dev Box.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"dev_center_attached_network_name": {
			Description:  "The name of the Dev Center Attached Network in parent Project of the Dev Center Project Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"managed_virtual_network_regions": {
			Description: "A list of the regions of the managed Virtual Network. When specified, the Dev Center Project Pool will use a Microsoft managed network.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     location.EnhancedValidate,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"single_sign_on_enabled": {
			Description: "Whether Dev Boxes in the Pool will have SSO enabled or disabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"stop_on_disconnect_grace_period_minutes": {
			Description:  "The specified time in minutes to wait before stopping a Dev Center Dev Box once disconnect is detected.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(60, 480),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
