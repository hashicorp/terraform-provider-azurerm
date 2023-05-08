package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorRouteRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorRouteName,
			},

			"cdn_frontdoor_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorEndpointID,
			},

			"cdn_frontdoor_origin_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			// I have to expose this here because that is the only way the user will
			// know that the resource has been modified outside of terraform...
			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeSet,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			// I have to expose this here because that is the only way the user will
			// know that the resource has been modified outside of terraform...
			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"cache": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"query_strings": {
							Type:     pluginsdk.TypeList,
							Computed: true,

							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"query_string_caching_behavior": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"compression_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"content_types_to_compress": {
							Type:     pluginsdk.TypeList,
							Computed: true,

							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"https_redirect_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"cdn_frontdoor_origin_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"patterns_to_match": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"cdn_frontdoor_rule_set_ids": {
				Type:     pluginsdk.TypeSet,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"supported_protocols": {
				Type:     pluginsdk.TypeSet,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceCdnFrontDoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	endpointIdRaw := d.Get("cdn_frontdoor_endpoint_id").(string)
	endpointId, err := parse.FrontDoorEndpointID(endpointIdRaw)
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorRouteID(subscriptionId, endpointId.ResourceGroup, endpointId.ProfileName, endpointId.AfdEndpointName, d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.RouteName)
	d.Set("cdn_frontdoor_endpoint_id", endpointId.ID())

	if props := resp.RouteProperties; props != nil {
		customDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
		if err != nil {
			return err
		}

		d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
		d.Set("enabled", flattenEnabledBool(props.EnabledState))
		d.Set("forwarding_protocol", props.ForwardingProtocol)
		d.Set("https_redirect_enabled", flattenHttpsRedirectToBool(props.HTTPSRedirect))
		d.Set("cdn_frontdoor_origin_path", props.OriginPath)
		d.Set("patterns_to_match", props.PatternsToMatch)
		d.Set("link_to_default_domain", flattenLinkToDefaultDomainToBool(props.LinkToDefaultDomain))

		if err := d.Set("cache", flattenCdnFrontdoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
			return fmt.Errorf("setting `cache`: %+v", err)
		}

		originGroupId, err := flattenOriginGroupResourceReference(props.OriginGroup)
		if err != nil {
			return fmt.Errorf("flattening `cdn_frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_origin_group_id", originGroupId); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_rule_set_ids", flattenRuleSetResourceArray(props.RuleSets)); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
		}

		if err := d.Set("supported_protocols", flattenCdnFrontdoorRouteEndpointProtocolsArray(props.SupportedProtocols)); err != nil {
			return fmt.Errorf("setting `supported_protocols`: %+v", err)
		}
	}

	return nil
}
