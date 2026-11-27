// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_network_acl

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func signalrServiceNetworkAclSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"signalr_service_id": {
			Description:  "The ID of the SignalR service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateSignalRID,
		},

		"default_action": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(signalr.PossibleValuesForACLAction(), false),
		},

		"public_network": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_request_types": {
						Description:   "The allowed request types for the public network.",
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"public_network.0.denied_request_types"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(signalr.PossibleValuesForSignalRRequestType(), false),
						},
					},

					"denied_request_types": {
						Description:   "The denied request types for the public network.",
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"public_network.0.allowed_request_types"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(signalr.PossibleValuesForSignalRRequestType(), false),
						},
					},
				},
			},
		},

		"private_endpoint": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description:  "The ID of the Private Endpoint which is based on the SignalR service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: privateendpoints.ValidatePrivateEndpointID,
					},

					"allowed_request_types": {
						Description: "The allowed request types for the Private Endpoint Connection.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(signalr.PossibleValuesForSignalRRequestType(), false),
						},
					},

					"denied_request_types": {
						Description: "The denied request types for the Private Endpoint Connection.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(signalr.PossibleValuesForSignalRRequestType(), false),
						},
					},
				},
			},
		},
	}
}
