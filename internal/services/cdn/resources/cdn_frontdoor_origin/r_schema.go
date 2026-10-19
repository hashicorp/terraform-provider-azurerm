// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorOriginName,
		},

		"cdn_frontdoor_origin_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: afdorigins.ValidateOriginGroupID,
		},

		"host_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_name_check_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"http_port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      80,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"https_port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      443,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"origin_host_header": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IsIPv6Address, validation.IsIPv4Address, validation.StringIsNotEmpty),
		},

		"priority": {
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "Access request for CDN FrontDoor Private Link Origin",
						ValidateFunc: validation.StringLenBetween(1, 140),
					},

					"target_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      500,
			ValidateFunc: validation.IntBetween(1, 1000),
		},
	}
}
