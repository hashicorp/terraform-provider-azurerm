// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnEndpointSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"profile_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"origin_host_header": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"is_http_allowed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"is_https_allowed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"origin": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"http_port": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
						Default:  80,
					},

					"https_port": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
						Default:  443,
					},
				},
			},
		},

		"origin_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"querystring_caching_behaviour": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(cdn.QueryStringCachingBehaviorIgnoreQueryString),
			ValidateFunc: validation.StringInEnumSlice(cdn.PossibleQueryStringCachingBehaviorValues(), false),
		},

		"content_types_to_compress": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"is_compression_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"probe_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"geo_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relative_path": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.ActionTypeAllow),
							string(cdn.ActionTypeBlock),
						}, false),
					},
					"country_codes": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"optimization_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInEnumSlice(cdn.PossibleOptimizationTypeValues(), false),
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"global_delivery_rule": endpointGlobalDeliveryRule(),

		"delivery_rule": helpers.EndpointDeliveryRule(),

		"tags": commonschema.Tags(),
	}
}
