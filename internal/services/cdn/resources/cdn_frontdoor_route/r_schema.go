// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorRouteSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorRouteName,
		},

		"cdn_frontdoor_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: routes.ValidateAfdEndpointID,
		},

		"cdn_frontdoor_origin_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: afdorigingroups.ValidateOriginGroupID,
		},

		// NOTE: These are not sent to the API, they are only here so Terraform
		// can provision/destroy the resources in the correct order.
		// Made this field optional to address comments in Issue #29063
		"cdn_frontdoor_origin_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: afdorigins.ValidateOriginGroupOriginID,
			},
		},

		"cdn_frontdoor_custom_domain_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: afddomains.ValidateCustomDomainID,
			},
		},

		"link_to_default_domain": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"cache": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"query_strings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringDoesNotContainAny(","),
						},
					},

					"query_string_caching_behavior": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(routes.AfdQueryStringCachingBehaviorIgnoreQueryString),
						ValidateFunc: validation.StringInSlice(routes.PossibleValuesForAfdQueryStringCachingBehavior(), false),
					},

					"compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"content_types_to_compress": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(helpers.FrontDoorContentTypes(), false),
						},
					},
				},
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"forwarding_protocol": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(routes.ForwardingProtocolMatchRequest),
			ValidateFunc: validation.StringInSlice(routes.PossibleValuesForForwardingProtocol(), false),
		},

		"https_redirect_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"cdn_frontdoor_origin_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"patterns_to_match": {
			Type:     pluginsdk.TypeList,
			Required: true,

			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cdn_frontdoor_rule_set_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: rulesets.ValidateRuleSetID,
			},
		},

		"supported_protocols": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MaxItems: 2,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(routes.PossibleValuesForAFDEndpointProtocols(), false),
			},
		},
	}
}
