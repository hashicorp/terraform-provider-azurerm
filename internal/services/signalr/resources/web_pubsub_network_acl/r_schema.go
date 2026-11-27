// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_network_acl

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func webPubsubNetworkAclSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(&webpubsub.WebPubSubId{}, "Web Pub/Sub Service"),

		"default_action": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webpubsub.ACLActionDeny,
			ValidateFunc: validation.StringInSlice(webpubsub.PossibleValuesForACLAction(), false),
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
							ValidateFunc: validation.StringInSlice(webpubsub.PossibleValuesForWebPubSubRequestType(), false),
						},
					},

					"denied_request_types": {
						Description:   "The denied request types for the public network.",
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"public_network.0.allowed_request_types"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(webpubsub.PossibleValuesForWebPubSubRequestType(), false),
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
						Description:  "The ID of the Private Endpoint which is based on the Web Pubsub service.",
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
							ValidateFunc: validation.StringInSlice(webpubsub.PossibleValuesForWebPubSubRequestType(), false),
						},
					},

					"denied_request_types": {
						Description: "The denied request types for the Private Endpoint Connection.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(webpubsub.PossibleValuesForWebPubSubRequestType(), false),
						},
					},
				},
			},
		},
	}
}
