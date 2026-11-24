// Copyright IBM Corp. 2014, 2026
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
			Description: "The name of the CDN Endpoint.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"profile_name": {
			Description: "The CDN Profile to which to attach the CDN Endpoint.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"origin_host_header": {
			Description: "The host header CDN provider will send along with content requests to origins.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"is_http_allowed": {
			Description: "Whether HTTP allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"is_https_allowed": {
			Description: "Whether HTTPS allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"origin": {
			Description: "The set of origins of the CDN endpoint. When multiple origins exist, the first origin will be used as primary and rest will be used as failover options.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the origin. This is an arbitrary value. However, this value needs to be unique under the endpoint.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
					},

					"host_name": {
						Description: "A string that determines the hostname/IP address of the origin server. This string can be a domain name, Storage Account endpoint, Web App endpoint, IPv4 address or IPv6 address.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
					},

					"http_port": {
						Description: "The HTTP port of the origin.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Default:     80,
					},

					"https_port": {
						Description: "The HTTPS port of the origin.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Default:     443,
					},
				},
			},
		},

		"origin_path": {
			Description: "The path used at for origin requests.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"querystring_caching_behaviour": {
			Description:  "Sets query string caching behaviour. `NotSet` value can be used for `Premium Verizon` CDN profile.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(cdn.QueryStringCachingBehaviorIgnoreQueryString),
			ValidateFunc: validation.StringInEnumSlice(cdn.PossibleQueryStringCachingBehaviorValues(), false),
		},

		"content_types_to_compress": {
			Description: "An array of strings that indicates a content types on which compression will be applied. The value for the elements should be MIME types.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"is_compression_enabled": {
			Description: "Whether compression is to be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"probe_path": {
			Description: "The path to a file hosted on the origin which helps accelerate delivery of the dynamic content and calculate the most optimal routes for the CDN. This is relative to the `origin_path`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"geo_filter": {
			Description: "A set of Geo Filters for this CDN Endpoint.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relative_path": {
						Description: "The relative path applicable to geo filter.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"action": {
						Description: "The Action of the Geo Filter.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.ActionTypeAllow),
							string(cdn.ActionTypeBlock),
						}, false),
					},
					"country_codes": {
						Description: "A List of two letter country codes (e.g. `US`, `GB`) to be associated with this Geo Filter.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"optimization_type": {
			Description:  "What types of optimization should this CDN Endpoint optimize for?",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInEnumSlice(cdn.PossibleOptimizationTypeValues(), false),
		},

		"fqdn": {
			Description: "The Fully Qualified Domain Name of the CDN Endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"global_delivery_rule": endpointGlobalDeliveryRule(),

		"delivery_rule": helpers.EndpointDeliveryRule(),

		"tags": commonschema.Tags(),
	}
}
