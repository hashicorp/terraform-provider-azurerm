// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorRouteCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	endpoint, err := routes.ParseAfdEndpointID(d.Get("cdn_frontdoor_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := routes.NewRouteID(endpoint.SubscriptionId, endpoint.ResourceGroupName, endpoint.ProfileName, endpoint.AfdEndpointName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_route", id.ID())
		}
	}

	// NOTE: If HTTPS Redirect is enabled the Supported Protocols must support both HTTP and HTTPS
	// This configuration does not cause an error when provisioned, however the http requests that
	// are supposed to be redirected to https remain http requests
	httpsRedirect := d.Get("https_redirect_enabled").(bool)
	protocols := d.Get("supported_protocols").(*pluginsdk.Set).List()
	if httpsRedirect {
		if err := validate.SupportsBothHttpAndHttps(protocols, "https_redirect_enabled"); err != nil {
			return err
		}
	}

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if !linkToDefaultDomain && len(customDomains) == 0 {
		return errors.New("`link_to_default_domain` cannot be disabled when no `cdn_frontdoor_custom_domain_ids` are defined")
	}

	if len(customDomains) != 0 {
		if err := helpers.ValidateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
			return err
		}
	}

	var originGroup *routes.ResourceReference
	if originGroupRaw := d.Get("cdn_frontdoor_origin_group_id").(string); originGroupRaw != "" {
		originGroup = expandCdnFrontDoorRouteResourceReference(originGroupRaw)
	}

	props := routes.Route{
		Properties: &routes.RouteProperties{
			CustomDomains:       helpers.ExpandCustomDomainActivatedResourceArray(customDomains),
			CacheConfiguration:  expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]any)),
			EnabledState:        expandCdnFrontDoorRouteEnabled(d.Get("enabled").(bool)),
			ForwardingProtocol:  pointer.ToEnum[routes.ForwardingProtocol](d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       expandCdnFrontDoorRouteHttpsRedirect(httpsRedirect),
			LinkToDefaultDomain: expandCdnFrontDoorRouteDefaultDomain(linkToDefaultDomain),
			OriginGroup:         originGroup,
			OriginPath:          pointer.ToOrNil(d.Get("cdn_frontdoor_origin_path").(string)),
			PatternsToMatch:     pluginsdk.ExpandStringSlice(d.Get("patterns_to_match").([]any)),
			RuleSets:            expandCdnFrontdoorRouteRuleSetReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List()),
			SupportedProtocols:  expandCdnFrontDoorRouteEndpointProtocolsArray(protocols),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if origins := d.Get("cdn_frontdoor_origin_ids").([]any); len(origins) != 0 {
		d.Set("cdn_frontdoor_origin_ids", origins)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}
