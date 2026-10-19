// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorRouteRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]any); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", originIds)
	}

	d.Set("name", id.RouteName)
	d.Set("cdn_frontdoor_endpoint_id", afdendpoints.NewAfdEndpointID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.AfdEndpointName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			customDomains, err := helpers.FlattenCdnFrontDoorRouteCustomDomainActivatedResourceArray(props.CustomDomains)
			if err != nil {
				return err
			}

			d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
			d.Set("enabled", pointer.From(props.EnabledState) == routes.EnabledStateEnabled)
			d.Set("forwarding_protocol", pointer.FromEnum(props.ForwardingProtocol))
			d.Set("https_redirect_enabled", pointer.From(props.HTTPSRedirect) == routes.HTTPSRedirectEnabled)
			d.Set("cdn_frontdoor_origin_path", props.OriginPath)
			d.Set("patterns_to_match", props.PatternsToMatch)
			d.Set("link_to_default_domain", pointer.From(props.LinkToDefaultDomain) == routes.LinkToDefaultDomainEnabled)

			if err := d.Set("cache", flattenCdnFrontDoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
				return fmt.Errorf("setting `cache`: %+v", err)
			}

			originGroupId, err := flattenCdnFrontDoorRouteOriginGroupResourceReference(props.OriginGroup)
			if err != nil {
				return fmt.Errorf("flattening `cdn_frontdoor_origin_group_id`: %+v", err)
			}

			if err := d.Set("cdn_frontdoor_origin_group_id", originGroupId); err != nil {
				return fmt.Errorf("setting `cdn_frontdoor_origin_group_id`: %+v", err)
			}

			ruleSetIDs, err := flattenCdnFrontDoorRouteRuleSetResourceArray(props.RuleSets)
			if err != nil {
				return fmt.Errorf("flattening `cdn_frontdoor_rule_set_ids`: %+v", err)
			}

			if err := d.Set("cdn_frontdoor_rule_set_ids", ruleSetIDs); err != nil {
				return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
			}

			if err := d.Set("supported_protocols", pointer.FromEnumSlice(props.SupportedProtocols)); err != nil {
				return fmt.Errorf("setting `supported_protocols`: %+v", err)
			}
		}
	}

	return nil
}
