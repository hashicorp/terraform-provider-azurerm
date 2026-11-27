// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this Network Manager Connectivity Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Description:  "The ID of the Network Manager.",
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
						Description:  "The group connectivity type.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForGroupConnectivity(), false),
					},

					"global_mesh_enabled": {
						Description: "Whether global mesh is supported for this group.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"network_group_id": {
						Description:  "The resource ID of Network Group which the configuration applies to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_hub_gateway": {
						Description: "Whether the hub gateway is used.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"connectivity_topology": {
			Description:  "The connectivity topology type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForConnectivityTopology(), false),
		},

		"connected_group_address_overlap_enabled": {
			Description: "Whether connected group address overlap is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"connected_group_private_endpoints_scale": {
			Description:  "The scale of Private Endpoints allowed in the connected group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(connectivityconfigurations.ConnectedGroupPrivateEndpointsScaleStandard),
			ValidateFunc: validation.StringInSlice(connectivityconfigurations.PossibleValuesForConnectedGroupPrivateEndpointsScale(), false),
		},

		"delete_existing_peering_enabled": {
			Description: "Whether to remove current existing Virtual Network Peering in the Connectivity Configuration affected scope.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"description": {
			Description:  "A description of the Connectivity Configuration.",
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
						Description:  "The resource ID used as hub in Hub And Spoke topology.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"resource_type": {
						Description:  "The resource Type used as hub in Hub And Spoke topology.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"global_mesh_enabled": {
			Description: "Whether global mesh is supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"peering_enforcement_enabled": {
			Description: "Whether peering enforcement is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
