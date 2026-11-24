// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorRouteSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Front Door Route.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorRouteName,
		},

		"cdn_frontdoor_endpoint_id": {
			Description:  "The resource ID of the Front Door Endpoint where this Front Door Route should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: routes.ValidateAfdEndpointID,
		},

		"cdn_frontdoor_origin_group_id": {
			Description:  "The resource ID of the Front Door Origin Group where this Front Door Route should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: afdorigingroups.ValidateOriginGroupID,
		},

		// NOTE: These are not sent to the API, they are only here so Terraform
		// can provision/destroy the resources in the correct order.
		// Made this field optional to address comments in Issue #29063
		"cdn_frontdoor_origin_ids": {
			Description: "One or more Front Door Origin resource IDs for this Front Door Route.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: afdorigins.ValidateOriginGroupOriginID,
			},
		},

		"cdn_frontdoor_custom_domain_ids": {
			Description: "The IDs of the Front Door Custom Domains which are associated with this Front Door Route.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: afddomains.ValidateCustomDomainID,
			},
		},

		"link_to_default_domain": {
			Description: "Should this Front Door Route be linked to the default endpoint?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"cache": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"query_strings": {
						Description: "Query strings to include or ignore.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringDoesNotContainAny(","),
						},
					},

					"query_string_caching_behavior": {
						Description:  "How the Front Door Route will cache requests that include query strings.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(routes.AfdQueryStringCachingBehaviorIgnoreQueryString),
						ValidateFunc: validation.StringInSlice(routes.PossibleValuesForAfdQueryStringCachingBehavior(), false),
					},

					"compression_enabled": {
						Description: "Is content compression enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"content_types_to_compress": {
						Description: "A list of one or more `Content types` (formerly known as `MIME types`) to compress.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(helpers.FrontDoorContentTypes(), false),
						},
					},
				},
			},
		},

		"enabled": {
			Description: "Is this Front Door Route enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"forwarding_protocol": {
			Description:  "The Protocol that will be use when forwarding traffic to backends.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(routes.ForwardingProtocolMatchRequest),
			ValidateFunc: validation.StringInSlice(routes.PossibleValuesForForwardingProtocol(), false),
		},

		"https_redirect_enabled": {
			Description: "Automatically redirect HTTP traffic to HTTPS traffic?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"cdn_frontdoor_origin_path": {
			Description: "A directory path on the Front Door Origin that can be used to retrieve content (e.g. `contoso.cloudapp.net/originpath`).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"patterns_to_match": {
			Description: "The route patterns of the rule.",
			Type:        pluginsdk.TypeList,
			Required:    true,

			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cdn_frontdoor_rule_set_ids": {
			Description: "A list of the Front Door Rule Set IDs which should be assigned to this Front Door Route.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: rulesets.ValidateRuleSetID,
			},
		},

		"supported_protocols": {
			Description: "One or more Protocols supported by this Front Door Route.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MaxItems:    2,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(routes.PossibleValuesForAFDEndpointProtocols(), false),
			},
		},
	}
}
