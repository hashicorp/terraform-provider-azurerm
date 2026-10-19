// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorRouteUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	// we need to lock the route for update because the custom domain
	// association may also be trying to update the route as well...
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

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

	if d.HasChange("cache") {
		props.CacheConfiguration = expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]any))
	}

	if d.HasChange("enabled") {
		props.EnabledState = expandCdnFrontDoorRouteEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("forwarding_protocol") {
		props.ForwardingProtocol = pointer.ToEnum[routes.ForwardingProtocol](d.Get("forwarding_protocol").(string))
	}

	if d.HasChange("https_redirect_enabled") {
		props.HTTPSRedirect = expandCdnFrontDoorRouteHttpsRedirect(httpsRedirect)
	}

	if d.HasChange("link_to_default_domain") {
		props.LinkToDefaultDomain = expandCdnFrontDoorRouteDefaultDomain(d.Get("link_to_default_domain").(bool))
	}

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		props.CustomDomains = helpers.ExpandCustomDomainActivatedResourceArray(customDomains)
	}

	if d.HasChange("cdn_frontdoor_origin_group_id") {
		props.OriginGroup = expandCdnFrontDoorRouteResourceReference(d.Get("cdn_frontdoor_origin_group_id").(string))
	}

	if d.HasChange("cdn_frontdoor_origin_path") {
		props.OriginPath = pointer.ToOrNil(d.Get("cdn_frontdoor_origin_path").(string))
	}

	if d.HasChange("patterns_to_match") {
		props.PatternsToMatch = pluginsdk.ExpandStringSlice(d.Get("patterns_to_match").([]any))
	}

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		props.RuleSets = expandCdnFrontdoorRouteRuleSetReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("supported_protocols") {
		props.SupportedProtocols = expandCdnFrontDoorRouteEndpointProtocolsArray(protocols)
	}

	if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]any); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", originIds)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}
