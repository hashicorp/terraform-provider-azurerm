// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_link_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateLinkServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Private Link Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"auto_approval_subscription_ids": {
			Description: "A list of Subscription UUID/GUID's that will be automatically be able to use this Private Link Service.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"destination_ip_address": {
			Description:  "The destination IP address of the Private Link Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
			ExactlyOneOf: []string{"load_balancer_frontend_ip_configuration_ids", "destination_ip_address"},
		},

		"proxy_protocol_enabled": {
			Description: "Whether the Private Link Service should support the proxy protocol.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"visibility_subscription_ids": {
			Description: "A list of Subscription UUID/GUID's that will be able to see this Private Link Service.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"*"}, false)),
			},
		},

		"fqdns": {
			Description: "List of FQDNs allowed for the Private Link Service.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		// Required by the API you can't create the resource without at least
		// one ip configuration once primary is set it is set forever unless
		// you destroy the resource and recreate it.
		"nat_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 8,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name which should be used for the NAT IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"private_ip_address": {
						Description:  "A Private Static IP Address for this IP Configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					// Only IPv4 is supported by the API, but I am exposing this
					// as they will support IPv6 in a future release.
					"private_ip_address_version": {
						Description: "The version of the IP Protocol which should be used. At this time the only.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(privatelinkservices.IPVersionIPvFour),
						ValidateFunc: validation.StringInSlice([]string{
							string(privatelinkservices.IPVersionIPvFour),
						}, false),
					},
					"subnet_id": {
						Description:  "The ID of the Subnet which should be used for the Private Link Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"primary": {
						Description: "Whether this is the primary IP configuration.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
				},
			},
		},

		"load_balancer_frontend_ip_configuration_ids": {
			Description:  "A list of Frontend IP Configuration IDs from a Standard Load Balancer, where traffic from the Private Link Service should be routed. You can use Load Balancer Rules to direct this traffic to appropriate backend pools where your applications are running.",
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			ForceNew:     true,
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"load_balancer_frontend_ip_configuration_ids", "destination_ip_address"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"alias": {
			Description: "A globally unique DNS Name for your Private Link Service. You can use this alias to request a connection to your Private Link Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
