// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_hardware_security_module

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dedicatedHardwareSecurityModuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DedicatedHardwareSecurityModuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dedicatedhsms.PossibleValuesForSkuName(), false),
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_interface_private_ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"management_network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_interface_private_ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"stamp_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"stamp1",
				"stamp2",
			}, false),
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
