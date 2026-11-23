// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppGateway.ApplicationGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgateways.ParseApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return resourceApplicationGatewaySetFlatten(d, id, resp.Model)
}

func base64EncodedStateFunc(v any) string {
	switch s := v.(type) {
	case string:
		return base64.EncodeIfNot(s)
	default:
		return ""
	}
}

func checkSslPolicy(sslPolicy []any) error {
	if len(sslPolicy) > 0 && sslPolicy[0] != nil {
		v := sslPolicy[0].(map[string]any)
		disabledProtocols := v["disabled_protocols"].([]any)
		policyType := v["policy_type"].(string)
		if len(disabledProtocols) > 0 && policyType != "" {
			return fmt.Errorf("setting disabled_protocols is not allowed when policy_type is defined")
		}
	}
	return nil
}

func checkBasicSkuFeatures(d *pluginsdk.ResourceDiff) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	if hasAutoscaleConfig {
		return fmt.Errorf("the Application Gateway does not support `autoscale_configuration` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	capacity, hasCapacityConfig := d.GetOk("sku.0.capacity")
	if hasCapacityConfig {
		if capacity.(int) > 2 || capacity.(int) < 1 {
			return fmt.Errorf("`capacity` value %q for the selected SKU tier %q is invalid. Value must be between [1-2]", capacity, applicationgateways.ApplicationGatewaySkuNameBasic)
		}
	} else {
		return fmt.Errorf("the Application Gateway must specify a `capacity` value between [1-2] for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	_, hasMtlsConfig := d.GetOk("trusted_client_certificate")
	if hasMtlsConfig {
		return fmt.Errorf("the Application Gateway does not support `trusted_client_certificate` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	rewriteRuleSet, hasRewriteRuleSetConfig := d.GetOk("rewrite_rule_set")
	if hasRewriteRuleSetConfig {
		for _, ruleSet := range rewriteRuleSet.([]any) {
			rs := ruleSet.(map[string]any)
			for _, rule := range rs["rewrite_rule"].([]any) {
				r := rule.(map[string]any)
				if len(r["url"].([]any)) > 0 {
					return fmt.Errorf("the Application Gateway does not support `url` inside the `rewrite_rule` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
				}
			}
		}
	}

	return nil
}

func applicationGatewayCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	capacity, hasCapacity := d.GetOk("sku.0.capacity")
	tier := d.Get("sku.0.tier").(string)

	if tier == string(applicationgateways.ApplicationGatewaySkuNameBasic) {
		if err := checkBasicSkuFeatures(d); err != nil {
			return err
		}
	} else if !hasAutoscaleConfig && !hasCapacity {
		return fmt.Errorf("the Application Gateway must specify either `capacity` or `autoscale_configuration` for the selected SKU tier %q", tier)
	}

	if err := checkSslPolicy(d.Get("ssl_policy").([]any)); err != nil {
		return err
	}

	sslProfiles := d.Get("ssl_profile").([]any)
	if len(sslProfiles) > 0 {
		for _, profile := range sslProfiles {
			if profile == nil {
				continue
			}
			v := profile.(map[string]any)
			if policy, ok := v["ssl_policy"]; ok && policy != nil {
				if err := checkSslPolicy(policy.([]any)); err != nil {
					return err
				}
			}

			if v["client_authentication_mode"].(string) == string(applicationgateways.ApplicationGatewayClientAuthVerificationModesPassthrough) {
				if v["verify_client_certificate_issuer_dn"].(bool) {
					return fmt.Errorf("`verify_client_certificate_issuer_dn` cannot be set to `true` when `client_authentication_mode` is `Passthrough` for `ssl_profile` %q", v["name"].(string))
				}
				if v["verify_client_certificate_revocation"].(string) == string(applicationgateways.ApplicationGatewayClientRevocationOptionsOCSP) {
					return fmt.Errorf("`verify_client_certificate_revocation` cannot be set to `OCSP` when `client_authentication_mode` is `Passthrough` for `ssl_profile` %q", v["name"].(string))
				}
			}
		}
	}

	if hasCapacity {
		if (strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierStandard)) || strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierWAF))) && (capacity.(int) < 1 || capacity.(int) > 32) {
			return fmt.Errorf("the value '%d' exceeds the maximum capacity allowed for a %q V1 SKU, the %q SKU must have a capacity value between 1 and 32", capacity, tier, tier)
		}

		if (strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierStandardVTwo)) || strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierWAFVTwo))) && (capacity.(int) < 1 || capacity.(int) > 125) {
			return fmt.Errorf("the value '%d' exceeds the maximum capacity allowed for a %q V2 SKU, the %q SKU must have a capacity value between 1 and 125", capacity, tier, tier)
		}
	}

	backendHttpSettings := d.Get("backend_http_settings").(*schema.Set).List()
	for _, rawSettings := range backendHttpSettings {
		settings := rawSettings.(map[string]any)

		if settings["sni_name"].(string) != "" && !settings["sni_validation_enabled"].(bool) {
			return fmt.Errorf("`sni_name` can only be set when `sni_validation_enabled` is set to `true` in `backend_http_settings` block `%s`", settings["name"].(string))
		}
	}

	if tier != "" && d.HasChange("sku.0.tier") && slices.Contains(validate.DeprecatedV1SkuTiers, tier) {
		return fmt.Errorf("new creation / update to `%s` SKU tier is no longer supported, please use supported SKU tiers: `Basic`, `Standard_v2`, `WAF_v2`, refer to https://aka.ms/V1retirement", tier)
	}

	if d.HasChange("sku.0.name") {
		skuName := d.Get("sku.0.name").(string)

		if skuName != "" && slices.Contains(validate.DeprecatedV1SkuNames, skuName) {
			return fmt.Errorf("new creation / update to `%s` SKU name is no longer supported, please use supported SKU names: `Basic`, `Standard_v2`, `WAF_v2`, refer to https://aka.ms/V1retirement", skuName)
		}
	}

	for _, raw := range d.Get("backend").([]any) {
		v := raw.(map[string]any)
		if v["host_name"].(string) != "" && strings.EqualFold(v["protocol"].(string), string(applicationgateways.ApplicationGatewayProtocolTcp)) {
			return fmt.Errorf("`host_name` cannot be set when `protocol` is `Tcp` for `backend` %q", v["name"].(string))
		}
	}

	for _, raw := range d.Get("backend_http_settings").(*pluginsdk.Set).List() {
		v := raw.(map[string]any)
		if v["host_name"].(string) != "" && v["pick_host_name_from_backend_address"].(bool) {
			return fmt.Errorf("only one of `host_name` or `pick_host_name_from_backend_address` can be set for `backend_http_settings` %q", v["name"].(string))
		}
	}

	for _, raw := range d.Get("listener").(*pluginsdk.Set).List() {
		v := raw.(map[string]any)
		if v["host_names"].(*pluginsdk.Set).Len() > 0 && strings.EqualFold(v["protocol"].(string), string(applicationgateways.ApplicationGatewayProtocolTcp)) {
			return fmt.Errorf("`host_names` cannot be set when `protocol` is `Tcp` for `listener` %q", v["name"].(string))
		}
		if strings.EqualFold(v["protocol"].(string), string(applicationgateways.ApplicationGatewayProtocolTls)) && v["ssl_certificate_name"].(string) == "" {
			return fmt.Errorf("`ssl_certificate_name` must be set when `protocol` is `Tls` for `listener` %q", v["name"].(string))
		}
	}

	for _, raw := range d.Get("probe").(*schema.Set).List() {
		v := raw.(map[string]any)
		protocol := v["protocol"].(string)
		isHTTPProtocol := strings.EqualFold(protocol, string(applicationgateways.ApplicationGatewayProtocolHTTP)) || strings.EqualFold(protocol, string(applicationgateways.ApplicationGatewayProtocolHTTPS))
		isTcpTlsProtocol := strings.EqualFold(protocol, string(applicationgateways.ApplicationGatewayProtocolTcp)) || strings.EqualFold(protocol, string(applicationgateways.ApplicationGatewayProtocolTls))

		hasHost := v["host"].(string) != ""
		hasPickHost := v["pick_host_name_from_backend_http_settings"].(bool)
		if isHTTPProtocol && (hasHost == hasPickHost) {
			return fmt.Errorf("exactly one of `host` or `pick_host_name_from_backend_http_settings` must be set when `protocol` is `Http` or `Https` for `probe` %q", v["name"].(string))
		}

		if isHTTPProtocol && v["path"].(string) == "" {
			return fmt.Errorf("`path` must be specified when `protocol` is `Http` or `Https` for `probe` %q", v["name"].(string))
		}

		if isTcpTlsProtocol {
			hasPath := v["path"].(string) != ""
			matchConfigs := v["match"].([]any)
			hasMatch := len(matchConfigs) > 0 && matchConfigs[0] != nil
			if hasHost || hasPickHost || hasPath || hasMatch {
				return fmt.Errorf("`host`, `pick_host_name_from_backend_http_settings`, `path`, and `match` cannot be set when `protocol` is `Tcp` or `Tls` for `probe` %q", v["name"].(string))
			}
		}

		if v["proxy_protocol_header_enabled"].(bool) && !isTcpTlsProtocol {
			return fmt.Errorf("`proxy_protocol_header_enabled` can only be set when `protocol` is `Tcp` or `Tls` for `probe` %q", v["name"].(string))
		}

		timeout := v["timeout"].(int)
		interval := v["interval"].(int)
		if timeout > interval {
			return fmt.Errorf("`timeout` must not be greater than `interval` for `probe` %q", v["name"].(string))
		}
	}

	return nil
}

func applicationGatewayHttpListnerHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["frontend_ip_configuration_name"].(string))
		buf.WriteString(m["frontend_port_name"].(string))
		buf.WriteString(m["protocol"].(string))
		if v, ok := m["host_name"]; ok {
			buf.WriteString(v.(string))
		}
		if hostNames, ok := m["host_names"]; ok {
			fmt.Fprintf(&buf, "%s-", hostNames.(*pluginsdk.Set).List())
		}
		if v, ok := m["ssl_certificate_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["require_sni"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if v, ok := m["firewall_policy_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
		if v, ok := m["ssl_profile_name"]; ok {
			buf.WriteString(v.(string))
		}
		if customErrorConfiguration, ok := m["custom_error_configuration"].([]any); ok {
			for _, customErrorAttrs := range customErrorConfiguration {
				if customErrorAttrs == nil {
					continue
				}
				customError := customErrorAttrs.(map[string]any)
				if statusCode, ok := customError["status_code"]; ok {
					buf.WriteString(statusCode.(string))
				}
				if pageUrl, ok := customError["custom_error_page_url"]; ok {
					buf.WriteString(pageUrl.(string))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayBackendSettingsHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		fmt.Fprintf(&buf, "%d", m["port"].(int))
		buf.WriteString(m["protocol"].(string))
		buf.WriteString(m["cookie_based_affinity"].(string))
		fmt.Fprintf(&buf, "%t", m["dedicated_backend_connection_enabled"].(bool))

		if v, ok := m["path"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["affinity_cookie_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["host_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["probe_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["pick_host_name_from_backend_address"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if v, ok := m["request_timeout"]; ok {
			fmt.Fprintf(&buf, "%d", v.(int))
		}
		if connectionDraining, ok := m["connection_draining"].([]any); ok {
			for _, ac := range connectionDraining {
				if ac == nil {
					continue
				}
				config := ac.(map[string]any)
				fmt.Fprintf(&buf, "%t", config["enabled"].(bool))
				fmt.Fprintf(&buf, "%d", config["drain_timeout_sec"].(int))
			}
		}
		if trustedRootCertificateNames, ok := m["trusted_root_certificate_names"]; ok {
			fmt.Fprintf(&buf, "%s", trustedRootCertificateNames.([]any))
		}
		if v, ok := m["certificate_chain_validation_enabled"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if v, ok := m["sni_validation_enabled"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if v, ok := m["sni_name"]; ok {
			buf.WriteString(v.(string))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewaySSLCertificate(v any) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))

		if v, ok := m["data"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["password"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["key_vault_secret_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayBackendAddressPool(v any) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))

		if fqdns, ok := m["fqdns"]; ok {
			fmt.Fprintf(&buf, "%s", fqdns.(*pluginsdk.Set).List())
		}
		if ips, ok := m["ip_addresses"]; ok {
			fmt.Fprintf(&buf, "%s", ips.(*pluginsdk.Set).List())
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayProbeHash(v any) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["protocol"].(string))
		buf.WriteString(m["path"].(string))
		fmt.Fprintf(&buf, "%d", m["interval"].(int))
		fmt.Fprintf(&buf, "%d", m["timeout"].(int))
		fmt.Fprintf(&buf, "%d", m["unhealthy_threshold"].(int))

		if v, ok := m["host"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["port"]; ok {
			fmt.Fprintf(&buf, "%d", v.(int))
		}
		if v, ok := m["pick_host_name_from_backend_http_settings"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if v, ok := m["minimum_servers"]; ok {
			fmt.Fprintf(&buf, "%d", v.(int))
		}
		if v, ok := m["proxy_protocol_header_enabled"]; ok {
			fmt.Fprintf(&buf, "%t", v.(bool))
		}
		if match, ok := m["match"]; ok {
			if attrs := match.([]any); len(attrs) == 1 && attrs[0] != nil {
				attr := attrs[0].(map[string]any)
				body := attr["body"].(string)
				statusCodes := attr["status_code"].([]any)

				// Only include in hash if it's not the default
				defaultMatch := body == "" && len(statusCodes) == 1 && statusCodes[0].(string) == "200-399"
				if !defaultMatch {
					fmt.Fprintf(&buf, "%s-%+v", body, statusCodes)
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayListenerHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["frontend_ip_configuration_name"].(string))
		buf.WriteString(m["frontend_port_name"].(string))
		buf.WriteString(m["protocol"].(string))

		if v, ok := m["ssl_certificate_name"]; ok {
			buf.WriteString(v.(string))
		}

		if v, ok := m["ssl_profile_name"]; ok {
			buf.WriteString(v.(string))
		}

		if hostNames, ok := m["host_names"]; ok {
			fmt.Fprintf(&buf, "%s-", hostNames.(*pluginsdk.Set).List())
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayRoutingRuleHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["backend_address_pool_name"].(string))
		buf.WriteString(m["backend_name"].(string))
		buf.WriteString(m["listener_name"].(string))
		fmt.Fprintf(&buf, "%d", m["priority"].(int))
	}

	return pluginsdk.HashString(buf.String())
}
