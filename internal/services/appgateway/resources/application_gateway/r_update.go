// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppGateway.ApplicationGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgateways.ParseApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if payload.Properties == nil {
		payload.Properties = &applicationgateways.ApplicationGatewayPropertiesFormat{}
	}

	if d.HasChange("http2_enabled") {
		payload.Properties.EnableHTTP2 = pointer.To(d.Get("http2_enabled").(bool))
	}

	if d.HasChange("trusted_root_certificate") {
		trustedRootCertificates, err := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]any))
		if err != nil {
			return fmt.Errorf("expanding `trusted_root_certificate`: %+v", err)
		}
		payload.Properties.TrustedRootCertificates = trustedRootCertificates
	}

	if d.HasChange("request_routing_rule") {
		requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `request_routing_rule`: %+v", err)
		}
		payload.Properties.RequestRoutingRules = requestRoutingRules
	}

	if d.HasChange("routing_rule") {
		payload.Properties.RoutingRules = expandApplicationGatewayRoutingRules(d.Get("routing_rule").(*pluginsdk.Set).List(), *id)
	}

	if d.HasChange("url_path_map") {
		urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `url_path_map`: %+v", err)
		}

		payload.Properties.UrlPathMaps = urlPathMaps
	}

	if d.HasChange("redirect_configuration") {
		redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `redirect_configuration`: %+v", err)
		}

		payload.Properties.RedirectConfigurations = redirectConfigurations
	}

	if d.HasChange("ssl_certificate") {
		sslCertificates, err := expandApplicationGatewaySslCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `ssl_certificate`: %+v", err)
		}

		payload.Properties.SslCertificates = sslCertificates
	}

	if d.HasChange("trusted_client_certificate") {
		trustedClientCertificates, err := expandApplicationGatewayTrustedClientCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `trusted_client_certificate`: %+v", err)
		}

		payload.Properties.TrustedClientCertificates = trustedClientCertificates
	}

	if d.HasChange("ssl_profile") {
		payload.Properties.SslProfiles = expandApplicationGatewaySslProfiles(d, id.ID())
	}

	gatewayIPConfigurations, stopApplicationGateway := expandApplicationGatewayIPConfigurations(d)
	if d.HasChange("gateway_ip_configuration") {
		payload.Properties.GatewayIPConfigurations = gatewayIPConfigurations
	}

	if d.HasChange("global") {
		payload.Properties.GlobalConfiguration = expandApplicationGatewayGlobalConfiguration(d.Get("global").([]any))
	}

	if d.HasChange("http_listener") {
		httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `http_listener`: %+v", err)
		}

		payload.Properties.HTTPListeners = httpListeners
	}

	if d.HasChange("listener") {
		payload.Properties.Listeners = expandApplicationGatewayListeners(d.Get("listener").(*schema.Set).List(), *id)
	}

	if d.HasChange("rewrite_rule_set") {
		rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
		if err != nil {
			return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
		}

		payload.Properties.RewriteRuleSets = rewriteRuleSets
	}

	if d.HasChange("autoscale_configuration") {
		payload.Properties.AutoscaleConfiguration = expandApplicationGatewayAutoscaleConfiguration(d)
	}

	if d.HasChange("custom_error_configuration") {
		payload.Properties.CustomErrorConfigurations = expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]any))
	}

	if d.HasChange("backend_address_pool") {
		payload.Properties.BackendAddressPools = expandApplicationGatewayBackendAddressPools(d)
	}

	if d.HasChange("backend_http_settings") {
		payload.Properties.BackendHTTPSettingsCollection = expandApplicationGatewayBackendHTTPSettings(d.Get("backend_http_settings").(*schema.Set).List(), id.ID())
	}

	if d.HasChange("backend") {
		payload.Properties.BackendSettingsCollection = expandApplicationGatewayBackendSettings(d.Get("backend").([]any), *id)
	}

	if d.HasChange("frontend_ip_configuration") {
		payload.Properties.FrontendIPConfigurations = expandApplicationGatewayFrontendIPConfigurations(d, id.ID())
	}

	if d.HasChange("frontend_port") {
		payload.Properties.FrontendPorts = expandApplicationGatewayFrontendPorts(d)
	}

	if d.HasChange("private_link_configuration") {
		payload.Properties.PrivateLinkConfigurations = expandApplicationGatewayPrivateLinkConfigurations(d)
	}

	if d.HasChange("probe") {
		payload.Properties.Probes = expandApplicationGatewayProbes(d.Get("probe").(*schema.Set).List())
	}

	if d.HasChange("sku") {
		payload.Properties.Sku = expandApplicationGatewaySku(d)
	}

	if d.HasChange("ssl_policy") {
		payload.Properties.SslPolicy = expandApplicationGatewaySslPolicy(d.Get("ssl_policy").([]any))
	}

	if d.HasChange("zones") {
		zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		if len(zones) > 0 {
			payload.Zones = &zones
		}
	}

	if d.HasChange("fips_enabled") {
		payload.Properties.EnableFips = pointer.To(d.Get("fips_enabled").(bool))
	}

	if d.HasChange("force_firewall_policy_association") {
		payload.Properties.ForceFirewallPolicyAssociation = pointer.To(d.Get("force_firewall_policy_association").(bool))
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		payload.Identity = expandedIdentity
	}

	if d.HasChange("waf_configuration") {
		payload.Properties.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(applicationgateways.ApplicationGatewayTierWAFVTwo) && wafFileUploadLimit > 500 {
		return fmt.Errorf("only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", applicationgateways.ApplicationGatewayTierWAFVTwo)
	}

	if d.HasChange("firewall_policy_id") {
		if d.Get("firewall_policy_id").(string) != "" {
			payload.Properties.FirewallPolicy = &applicationgateways.SubResource{
				Id: pointer.To(d.Get("firewall_policy_id").(string)),
			}
		} else {
			payload.Properties.FirewallPolicy = nil
		}
	}

	if stopApplicationGateway {
		if err := client.StopThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if stopApplicationGateway {
		if err := client.StartThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceApplicationGatewayRead(d, meta)
}
