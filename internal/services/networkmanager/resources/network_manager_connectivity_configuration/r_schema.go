// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.HasChange("connected_group_address_overlap_enabled") {
				oldVal, newVal := rd.GetChange("connected_group_address_overlap_enabled")

				if oldVal.(bool) && !newVal.(bool) {
					if err := rd.ForceNew("connected_group_address_overlap_enabled"); err != nil {
						return err
					}
				}
			}

			topology := rd.Get("connectivity_topology").(string)
			peeringEnforcementEnabled := rd.Get("peering_enforcement_enabled").(bool)

			if topology == string(connectivityconfigurations.ConnectivityTopologyMesh) &&
				peeringEnforcementEnabled {
				return errors.New("`peering_enforcement_enabled` can only be set to `true` when `connectivity_topology` is `HubAndSpoke`")
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: connectivityconfigurations.ValidateNetworkManagerID,
		},

		"applies_to_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_connectivity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForGroupConnectivity(), false),
					},

					"global_mesh_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"network_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_hub_gateway": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"connectivity_topology": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForConnectivityTopology(), false),
		},

		"connected_group_address_overlap_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"connected_group_private_endpoints_scale": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(connectivityconfigurations.ConnectedGroupPrivateEndpointsScaleStandard),
			ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForConnectedGroupPrivateEndpointsScale(), false),
		},

		"delete_existing_peering_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"resource_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"global_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"peering_enforcement_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
