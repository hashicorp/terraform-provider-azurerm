// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorOriginSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Front Door Origin.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorOriginName,
		},

		"cdn_frontdoor_origin_group_id": {
			Description:  "The ID of the Front Door Origin Group within which this Front Door Origin should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: afdorigins.ValidateOriginGroupID,
		},

		"host_name": {
			Description:  "The IPv4 address, IPv6 address or Domain name of the Origin.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_name_check_enabled": {
			Description: "Whether certificate name checks are enabled for this origin.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"enabled": {
			Description: "Should the origin be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"http_port": {
			Description:  "The value of the HTTP port.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      80,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"https_port": {
			Description:  "The value of the HTTPS port.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      443,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"origin_host_header": {
			Description:  "The host header value (an IPv4 address, IPv6 address or Domain name) which is sent to the origin with each request. If unspecified the hostname from the request will be used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IsIPv6Address, validation.IsIPv4Address, validation.StringIsNotEmpty),
		},

		"priority": {
			Description:  "Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 5),
		},

		"private_link": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.Location(),

					"private_link_target_id": {
						Description: "The ID of the Azure Resource to connect to via the Private Link.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.Any(
							apimanagementservice.ValidateServiceID,
							applicationgateways.ValidateApplicationGatewayID,
							commonids.ValidateAppServiceID,
							commonids.ValidateStorageAccountID,
							managedenvironments.ValidateManagedEnvironmentID,
							privatelinkservices.ValidatePrivateLinkServiceID,
						),
					},

					"request_message": {
						Description:  "The request message that will be submitted to the `private_link_target_id` when requesting the Private Link endpoint connection.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "Access request for CDN FrontDoor Private Link Origin",
						ValidateFunc: validation.StringLenBetween(1, 140),
					},

					"target_type": {
						Description: "The type of target for this Private Link Endpoint.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"blob",
							"blob_secondary",
							"Gateway",
							"managedEnvironments",
							"sites",
							"web",
							"web_secondary",
						}, false),
					},
				},
			},
		},

		"weight": {
			Description:  "The weight of the origin in a given origin group for load balancing.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      500,
			ValidateFunc: validation.IntBetween(1, 1000),
		},
	}
}
