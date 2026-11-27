// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppGateway.ApplicationGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationgateways.NewApplicationGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_application_gateway", id.ID())
		}
	}

	http2Enabled := d.Get("http2_enabled").(bool)
	t := d.Get("tags").(map[string]any)

	// Gateway ID is needed to link sub-resources together in expand functions
	trustedRootCertificates, err := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]any))
	if err != nil {
		return fmt.Errorf("expanding `trusted_root_certificate`: %+v", err)
	}

	requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `request_routing_rule`: %+v", err)
	}

	urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `url_path_map`: %+v", err)
	}

	redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `redirect_configuration`: %+v", err)
	}

	sslCertificates, err := expandApplicationGatewaySslCertificates(d)
	if err != nil {
		return fmt.Errorf("expanding `ssl_certificate`: %+v", err)
	}

	trustedClientCertificates, err := expandApplicationGatewayTrustedClientCertificates(d)
	if err != nil {
		return fmt.Errorf("expanding `trusted_client_certificate`: %+v", err)
	}

	sslProfiles := expandApplicationGatewaySslProfiles(d, id.ID())

	gatewayIPConfigurations, _ := expandApplicationGatewayIPConfigurations(d)

	globalConfiguration := expandApplicationGatewayGlobalConfiguration(d.Get("global").([]any))

	httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `http_listener`: %+v", err)
	}

	rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
	if err != nil {
		return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
	}

	gateway := applicationgateways.ApplicationGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(t),
		Properties: &applicationgateways.ApplicationGatewayPropertiesFormat{
			AutoscaleConfiguration:        expandApplicationGatewayAutoscaleConfiguration(d),
			TrustedRootCertificates:       trustedRootCertificates,
			CustomErrorConfigurations:     expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]any)),
			BackendAddressPools:           expandApplicationGatewayBackendAddressPools(d),
			BackendHTTPSettingsCollection: expandApplicationGatewayBackendHTTPSettings(d.Get("backend_http_settings").(*schema.Set).List(), id.ID()),
			BackendSettingsCollection:     expandApplicationGatewayBackendSettings(d.Get("backend").([]any), id),
			EnableHTTP2:                   pointer.To(http2Enabled),
			FrontendIPConfigurations:      expandApplicationGatewayFrontendIPConfigurations(d, id.ID()),
			FrontendPorts:                 expandApplicationGatewayFrontendPorts(d),
			GatewayIPConfigurations:       gatewayIPConfigurations,
			GlobalConfiguration:           globalConfiguration,
			HTTPListeners:                 httpListeners,
			Listeners:                     expandApplicationGatewayListeners(d.Get("listener").(*schema.Set).List(), id),
			PrivateLinkConfigurations:     expandApplicationGatewayPrivateLinkConfigurations(d),
			Probes:                        expandApplicationGatewayProbes(d.Get("probe").(*schema.Set).List()),
			RequestRoutingRules:           requestRoutingRules,
			RoutingRules:                  expandApplicationGatewayRoutingRules(d.Get("routing_rule").(*pluginsdk.Set).List(), id),
			RedirectConfigurations:        redirectConfigurations,
			Sku:                           expandApplicationGatewaySku(d),
			SslCertificates:               sslCertificates,
			TrustedClientCertificates:     trustedClientCertificates,
			SslProfiles:                   sslProfiles,
			SslPolicy:                     expandApplicationGatewaySslPolicy(d.Get("ssl_policy").([]any)),

			RewriteRuleSets: rewriteRuleSets,
			UrlPathMaps:     urlPathMaps,
		},
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		gateway.Zones = &zones
	}

	if v, ok := d.GetOk("fips_enabled"); ok {
		gateway.Properties.EnableFips = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("force_firewall_policy_association"); ok {
		gateway.Properties.ForceFirewallPolicyAssociation = pointer.To(v.(bool))
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		gateway.Identity = expandedIdentity
	}

	if _, ok := d.GetOk("waf_configuration"); ok {
		gateway.Properties.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(applicationgateways.ApplicationGatewayTierWAFVTwo) && wafFileUploadLimit > 500 {
		return fmt.Errorf("only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", applicationgateways.ApplicationGatewayTierWAFVTwo)
	}

	if v, ok := d.GetOk("firewall_policy_id"); ok {
		gateway.Properties.FirewallPolicy = &applicationgateways.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, gateway, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceApplicationGatewayRead(d, meta)
}
