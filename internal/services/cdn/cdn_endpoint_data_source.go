// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnEndpointRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"origin_host_header": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"is_http_allowed": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_https_allowed": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"origin": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"http_port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"https_port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"origin_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"querystring_caching_behaviour": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"content_types_to_compress": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"is_compression_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"probe_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"geo_filter": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"relative_path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"country_codes": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"optimization_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceCdnEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("profile_name", id.ProfileName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.EndpointProperties; props != nil {
		d.Set("fqdn", props.HostName)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)
		d.Set("probe_path", props.ProbePath)
		d.Set("optimization_type", string(props.OptimizationType))

		compressionEnabled := false
		if v := props.IsCompressionEnabled; v != nil {
			compressionEnabled = *v
		}
		d.Set("is_compression_enabled", compressionEnabled)

		contentTypes := flattenAzureRMCdnEndpointContentTypes(props.ContentTypesToCompress)
		if err := d.Set("content_types_to_compress", contentTypes); err != nil {
			return fmt.Errorf("setting `content_types_to_compress`: %+v", err)
		}

		geoFilters := flattenCdnEndpointGeoFilters(props.GeoFilters)
		if err := d.Set("geo_filter", geoFilters); err != nil {
			return fmt.Errorf("setting `geo_filter`: %+v", err)
		}

		origins := flattenAzureRMCdnEndpointOrigin(props.Origins)
		if err := d.Set("origin", origins); err != nil {
			return fmt.Errorf("setting `origin`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
