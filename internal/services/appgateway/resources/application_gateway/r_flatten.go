// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceApplicationGatewaySetFlatten(d *pluginsdk.ResourceData, id *applicationgateways.ApplicationGatewayId, model *applicationgateways.ApplicationGateway) error {
	d.Set("name", id.ApplicationGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		if err = d.Set("identity", identity); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			if err = d.Set("trusted_root_certificate", flattenApplicationGatewayTrustedRootCertificates(props.TrustedRootCertificates, d)); err != nil {
				return fmt.Errorf("setting `trusted_root_certificate`: %+v", err)
			}

			if setErr := d.Set("backend_address_pool", flattenApplicationGatewayBackendAddressPools(props.BackendAddressPools)); setErr != nil {
				return fmt.Errorf("setting `backend_address_pool`: %+v", setErr)
			}

			backendHttpSettings, err := flattenApplicationGatewayBackendHTTPSettings(props.BackendHTTPSettingsCollection)
			if err != nil {
				return fmt.Errorf("flattening `backend_http_settings`: %+v", err)
			}

			if !features.SixPointOh() {
				backendHttpSettings = backendHTTPSettingsUnderUsedNames(d, backendHttpSettings)
			}

			if setErr := d.Set("backend_http_settings", backendHttpSettings); setErr != nil {
				return fmt.Errorf("setting `backend_http_settings`: %+v", setErr)
			}

			backendSettings, err := flattenApplicationGatewayBackendSettings(props.BackendSettingsCollection)
			if err != nil {
				return fmt.Errorf("flattening `backend`: %+v", err)
			}
			if setErr := d.Set("backend", backendSettings); setErr != nil {
				return fmt.Errorf("setting `backend`: %+v", setErr)
			}

			if setErr := d.Set("ssl_policy", flattenApplicationGatewaySslPolicy(props.SslPolicy)); setErr != nil {
				return fmt.Errorf("setting `ssl_policy`: %+v", setErr)
			}

			d.Set("http2_enabled", props.EnableHTTP2)
			d.Set("fips_enabled", props.EnableFips)
			d.Set("force_firewall_policy_association", props.ForceFirewallPolicyAssociation)

			httpListeners, err := flattenApplicationGatewayHTTPListeners(props.HTTPListeners)
			if err != nil {
				return fmt.Errorf("flattening `http_listener`: %+v", err)
			}
			if setErr := d.Set("http_listener", httpListeners); setErr != nil {
				return fmt.Errorf("setting `http_listener`: %+v", setErr)
			}

			listeners, err := flattenApplicationGatewayListeners(props.Listeners)
			if err != nil {
				return fmt.Errorf("flattening `listener`: %+v", err)
			}
			if setErr := d.Set("listener", listeners); setErr != nil {
				return fmt.Errorf("setting `listener`: %+v", setErr)
			}

			if setErr := d.Set("frontend_port", flattenApplicationGatewayFrontendPorts(props.FrontendPorts)); setErr != nil {
				return fmt.Errorf("setting `frontend_port`: %+v", setErr)
			}

			frontendIPConfigurations, err := flattenApplicationGatewayFrontendIPConfigurations(props.FrontendIPConfigurations)
			if err != nil {
				return fmt.Errorf("flattening `frontend IP configuration`: %+v", err)
			}
			if setErr := d.Set("frontend_ip_configuration", frontendIPConfigurations); setErr != nil {
				return fmt.Errorf("setting `frontend_ip_configuration`: %+v", setErr)
			}

			if setErr := d.Set("gateway_ip_configuration", flattenApplicationGatewayIPConfigurations(props.GatewayIPConfigurations)); setErr != nil {
				return fmt.Errorf("setting `gateway_ip_configuration`: %+v", setErr)
			}

			if setErr := d.Set("global", flattenApplicationGatewayGlobalConfiguration(props.GlobalConfiguration)); setErr != nil {
				return fmt.Errorf("setting `global`: %+v", setErr)
			}

			if setErr := d.Set("private_endpoint_connection", flattenApplicationGatewayPrivateEndpoints(props.PrivateEndpointConnections)); setErr != nil {
				return fmt.Errorf("setting `private_endpoint_connection`: %+v", setErr)
			}

			if setErr := d.Set("private_link_configuration", flattenApplicationGatewayPrivateLinkConfigurations(props.PrivateLinkConfigurations)); setErr != nil {
				return fmt.Errorf("setting `private_link_configuration`: %+v", setErr)
			}

			if setErr := d.Set("probe", flattenApplicationGatewayProbes(props.Probes)); setErr != nil {
				return fmt.Errorf("setting `probe`: %+v", setErr)
			}

			requestRoutingRules, err := flattenApplicationGatewayRequestRoutingRules(props.RequestRoutingRules)
			if err != nil {
				return fmt.Errorf("flattening `request_routing_rule`: %+v", err)
			}
			if setErr := d.Set("request_routing_rule", requestRoutingRules); setErr != nil {
				return fmt.Errorf("setting `request_routing_rule`: %+v", setErr)
			}

			routingRules, err := flattenApplicationGatewayRoutingRules(props.RoutingRules)
			if err != nil {
				return fmt.Errorf("flattening `routing_rule`: %+v", err)
			}
			if setErr := d.Set("routing_rule", routingRules); setErr != nil {
				return fmt.Errorf("setting `routing_rule`: %+v", setErr)
			}

			redirectConfigurations, err := flattenApplicationGatewayRedirectConfigurations(props.RedirectConfigurations)
			if err != nil {
				return fmt.Errorf("flattening `redirect configuration`: %+v", err)
			}
			if setErr := d.Set("redirect_configuration", redirectConfigurations); setErr != nil {
				return fmt.Errorf("setting `redirect_configuration`: %+v", setErr)
			}

			if setErr := d.Set("rewrite_rule_set", flattenApplicationGatewayRewriteRuleSets(props.RewriteRuleSets)); setErr != nil {
				return fmt.Errorf("setting `rewrite_rule_set`: %+v", setErr)
			}

			if setErr := d.Set("sku", flattenApplicationGatewaySku(props.Sku)); setErr != nil {
				return fmt.Errorf("setting `sku`: %+v", setErr)
			}

			if setErr := d.Set("autoscale_configuration", flattenApplicationGatewayAutoscaleConfiguration(props.AutoscaleConfiguration)); setErr != nil {
				return fmt.Errorf("setting `autoscale_configuration`: %+v", setErr)
			}

			if setErr := d.Set("ssl_certificate", flattenApplicationGatewaySslCertificates(props.SslCertificates, d)); setErr != nil {
				return fmt.Errorf("setting `ssl_certificate`: %+v", setErr)
			}

			if setErr := d.Set("trusted_client_certificate", flattenApplicationGatewayTrustedClientCertificates(props.TrustedClientCertificates)); setErr != nil {
				return fmt.Errorf("setting `trusted_client_certificate`: %+v", setErr)
			}

			sslProfiles, err := flattenApplicationGatewaySslProfiles(props.SslProfiles)
			if err != nil {
				return fmt.Errorf("flattening `ssl_profile`: %+v", err)
			}
			if setErr := d.Set("ssl_profile", sslProfiles); setErr != nil {
				return fmt.Errorf("setting `ssl_profile`: %+v", setErr)
			}

			if setErr := d.Set("custom_error_configuration", flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)); setErr != nil {
				return fmt.Errorf("setting `custom_error_configuration`: %+v", setErr)
			}

			urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.UrlPathMaps)
			if err != nil {
				return fmt.Errorf("flattening `url_path_map`: %+v", err)
			}
			if setErr := d.Set("url_path_map", urlPathMaps); setErr != nil {
				return fmt.Errorf("setting `url_path_map`: %+v", setErr)
			}

			if setErr := d.Set("waf_configuration", flattenApplicationGatewayWafConfig(props.WebApplicationFirewallConfiguration)); setErr != nil {
				return fmt.Errorf("setting `waf_configuration`: %+v", setErr)
			}

			firewallPolicyId := ""
			if props.FirewallPolicy != nil && props.FirewallPolicy.Id != nil {
				firewallPolicyId = *props.FirewallPolicy.Id
				policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(firewallPolicyId)
				if err == nil {
					firewallPolicyId = policyId.ID()
				}
			}
			d.Set("firewall_policy_id", firewallPolicyId)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenApplicationGatewayTrustedRootCertificates(certs *[]applicationgateways.ApplicationGatewayTrustedRootCertificate, d *pluginsdk.ResourceData) []any {
	results := make([]any, 0)
	if certs == nil {
		return results
	}

	// since the certificate data isn't returned lets load any existing data
	nameToDataMap := map[string]string{}
	if existing, ok := d.GetOk("trusted_root_certificate"); ok && existing != nil {
		for _, c := range existing.([]any) {
			b := c.(map[string]any)
			nameToDataMap[b["name"].(string)] = b["data"].(string)
		}
	}

	for _, cert := range *certs {
		output := map[string]any{}

		if v := cert.Id; v != nil {
			output["id"] = *v
		}

		kvsid := ""
		if props := cert.Properties; props != nil {
			if v := props.KeyVaultSecretId; v != nil {
				kvsid = *v
			}
		}
		output["key_vault_secret_id"] = kvsid

		if v := cert.Name; v != nil {
			output["name"] = *v

			// if there's no key vauld ID and we have a name, so try and look up the old data to pass it along
			if data, ok := nameToDataMap[*v]; ok && data != "" {
				output["data"] = data
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayBackendAddressPools(input *[]applicationgateways.ApplicationGatewayBackendAddressPool) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		ipAddressList := make([]any, 0)
		fqdnList := make([]any, 0)

		if props := config.Properties; props != nil {
			if props.BackendAddresses != nil {
				for _, address := range *props.BackendAddresses {
					if address.IPAddress != nil {
						ipAddressList = append(ipAddressList, *address.IPAddress)
					} else if address.Fqdn != nil {
						fqdnList = append(fqdnList, *address.Fqdn)
					}
				}
			}
		}

		output := map[string]any{
			"fqdns":        fqdnList,
			"ip_addresses": ipAddressList,
		}

		if config.Id != nil {
			output["id"] = *config.Id
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayBackendHTTPSettings(input *[]applicationgateways.ApplicationGatewayBackendHTTPSettings) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			output["cookie_based_affinity_enabled"] = pointer.From(props.CookieBasedAffinity) == applicationgateways.ApplicationGatewayCookieBasedAffinityEnabled
			if !features.SixPointOh() {
				output["cookie_based_affinity"] = string(pointer.From(props.CookieBasedAffinity))
			}

			output["dedicated_backend_connection_enabled"] = pointer.From(props.DedicatedBackendConnection)

			if affinityCookieName := props.AffinityCookieName; affinityCookieName != nil {
				output["affinity_cookie_name"] = affinityCookieName
			}

			if path := props.Path; path != nil {
				output["path"] = *path
			}

			certificateChainValidationEnabled := true
			if props.ValidateCertChainAndExpiry != nil {
				certificateChainValidationEnabled = *props.ValidateCertChainAndExpiry
			}
			output["certificate_chain_validation_enabled"] = certificateChainValidationEnabled

			output["connection_draining"] = flattenApplicationGatewayConnectionDraining(props.ConnectionDraining)

			if port := props.Port; port != nil {
				output["port"] = int(*port)
			}

			if hostName := props.HostName; hostName != nil {
				output["host_name"] = *hostName
			}

			if pickHostNameFromBackendAddress := props.PickHostNameFromBackendAddress; pickHostNameFromBackendAddress != nil {
				output["pick_host_name_from_backend_address"] = *pickHostNameFromBackendAddress
			}

			output["protocol"] = props.Protocol

			if timeout := props.RequestTimeout; timeout != nil {
				output["request_timeout"] = int(*timeout)
			}

			output["sni_name"] = pointer.From(props.SniName)

			sniValidationEnabled := true
			if props.ValidateSNI != nil {
				sniValidationEnabled = *props.ValidateSNI
			}
			output["sni_validation_enabled"] = sniValidationEnabled

			trustedRootCertificateNames := make([]any, 0)
			if certs := props.TrustedRootCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.TrustedRootCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					trustedRootCertificateNames = append(trustedRootCertificateNames, certId.Name)
				}
			}
			output["trusted_root_certificate_names"] = trustedRootCertificateNames

			if probe := props.Probe; probe != nil {
				if probe.Id != nil {
					id, err := parse.ProbeIDInsensitively(*probe.Id)
					if err != nil {
						return results, err
					}

					output["probe_name"] = id.Name
					output["probe_id"] = id.ID()
				}
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayConnectionDraining(input *applicationgateways.ApplicationGatewayConnectionDraining) []any {
	if input == nil {
		return []any{}
	}

	return []any{map[string]any{
		"enabled":           input.Enabled,
		"drain_timeout_sec": input.DrainTimeoutInSec,
	}}
}

func flattenApplicationGatewayBackendSettings(input *[]applicationgateways.ApplicationGatewayBackendSettings) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{
			"id":   pointer.From(v.Id),
			"name": pointer.From(v.Name),
		}

		if props := v.Properties; props != nil {
			output["port"] = int(pointer.From(props.Port))
			output["host_name"] = pointer.From(props.HostName)
			output["protocol"] = pointer.FromEnum(props.Protocol)
			output["timeout_in_seconds"] = int(pointer.From(props.Timeout))
			output["client_ip_preservation_enabled"] = pointer.From(props.EnableL4ClientIPPreservation)

			trustedRootCertificateNames := make([]any, 0)
			if certs := props.TrustedRootCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.TrustedRootCertificateIDInsensitively(pointer.From(cert.Id))
					if err != nil {
						return nil, err
					}

					trustedRootCertificateNames = append(trustedRootCertificateNames, certId.Name)
				}
			}
			output["trusted_root_certificate_names"] = trustedRootCertificateNames

			if probe := props.Probe; probe != nil {
				if probe.Id == nil {
					continue
				}

				id, err := parse.ProbeIDInsensitively(*probe.Id)
				if err != nil {
					return results, err
				}

				output["probe_name"] = id.Name
				output["probe_id"] = id.ID()
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewaySslPolicy(input *applicationgateways.ApplicationGatewaySslPolicy) []any {
	results := make([]any, 0)

	if input == nil {
		return results
	}

	output := map[string]any{}
	output["policy_name"] = input.PolicyName
	output["policy_type"] = input.PolicyType
	output["min_protocol_version"] = input.MinProtocolVersion

	cipherSuites := make([]any, 0)
	if input.CipherSuites != nil {
		for _, v := range *input.CipherSuites {
			cipherSuites = append(cipherSuites, string(v))
		}
	}
	output["cipher_suites"] = cipherSuites

	disabledSslProtocols := make([]any, 0)
	if input.DisabledSslProtocols != nil {
		for _, v := range *input.DisabledSslProtocols {
			disabledSslProtocols = append(disabledSslProtocols, string(v))
		}
	}
	output["disabled_protocols"] = disabledSslProtocols

	results = append(results, output)
	return results
}

func flattenApplicationGatewayHTTPListeners(input *[]applicationgateways.ApplicationGatewayHTTPListener) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if port := props.FrontendPort; port != nil {
				if port.Id != nil {
					portId, err := parse.FrontendPortIDInsensitively(*port.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_port_name"] = portId.Name
					output["frontend_port_id"] = portId.ID()
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.Id != nil {
					feConfigId, err := parse.FrontendIPConfigurationIDInsensitively(*feConfig.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_ip_configuration_name"] = feConfigId.Name
					output["frontend_ip_configuration_id"] = feConfigId.ID()
				}
			}

			if hostname := props.HostName; hostname != nil {
				output["host_name"] = *hostname
			}

			if hostnames := props.HostNames; hostnames != nil {
				output["host_names"] = pluginsdk.FlattenSlice(hostnames)
			}

			output["protocol"] = props.Protocol

			if cert := props.SslCertificate; cert != nil {
				if cert.Id != nil {
					certId, err := parse.SslCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_certificate_name"] = certId.Name
					output["ssl_certificate_id"] = certId.ID()
				}
			}

			if sni := props.RequireServerNameIndication; sni != nil {
				output["require_sni"] = *sni
			}

			if fwp := props.FirewallPolicy; fwp != nil && fwp.Id != nil {
				policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.Id)
				if err != nil {
					return nil, err
				}
				output["firewall_policy_id"] = policyId.ID()
			}

			if sslp := props.SslProfile; sslp != nil {
				if sslp.Id != nil {
					sslProfileId, err := parse.SslProfileIDInsensitively(*sslp.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_profile_name"] = sslProfileId.Name
					output["ssl_profile_id"] = sslProfileId.ID()
				}
			}

			output["custom_error_configuration"] = flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayListeners(input *[]applicationgateways.ApplicationGatewayListener) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{
			"id":   pointer.From(v.Id),
			"name": pointer.From(v.Name),
		}

		if props := v.Properties; props != nil {
			output["protocol"] = pointer.From(props.Protocol)
			output["host_names"] = pointer.From(props.HostNames)

			if port := props.FrontendPort; port != nil {
				if port.Id != nil {
					portId, err := parse.FrontendPortIDInsensitively(*port.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_port_name"] = portId.Name
					output["frontend_port_id"] = portId.ID()
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.Id != nil {
					feConfigId, err := parse.FrontendIPConfigurationIDInsensitively(*feConfig.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_ip_configuration_name"] = feConfigId.Name
					output["frontend_ip_configuration_id"] = feConfigId.ID()
				}
			}

			if cert := props.SslCertificate; cert != nil {
				if cert.Id != nil {
					certId, err := parse.SslCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_certificate_name"] = certId.Name
					output["ssl_certificate_id"] = certId.ID()
				}
			}

			if sslp := props.SslProfile; sslp != nil {
				if sslp.Id != nil {
					sslProfileId, err := parse.SslProfileIDInsensitively(*sslp.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_profile_name"] = sslProfileId.Name
					output["ssl_profile_id"] = sslProfileId.ID()
				}
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayIPConfigurations(input *[]applicationgateways.ApplicationGatewayIPConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if subnet := props.Subnet; subnet != nil {
				if subnet.Id != nil {
					output["subnet_id"] = *subnet.Id
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayGlobalConfiguration(input *applicationgateways.ApplicationGatewayGlobalConfiguration) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)

	if input.EnableRequestBuffering != nil {
		output["request_buffering_enabled"] = *input.EnableRequestBuffering
	}

	if input.EnableResponseBuffering != nil {
		output["response_buffering_enabled"] = *input.EnableResponseBuffering
	}

	return []any{output}
}

func flattenApplicationGatewayFrontendPorts(input *[]applicationgateways.ApplicationGatewayFrontendPort) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if props.Port != nil {
				output["port"] = int(*props.Port)
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayFrontendIPConfigurations(input *[]applicationgateways.ApplicationGatewayFrontendIPConfiguration) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		output := make(map[string]any)
		if config.Id != nil {
			output["id"] = *config.Id
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		if props := config.Properties; props != nil {
			output["private_ip_address_allocation"] = pointer.FromEnum(props.PrivateIPAllocationMethod)

			if props.Subnet != nil && props.Subnet.Id != nil {
				output["subnet_id"] = *props.Subnet.Id
			}

			if props.PrivateIPAddress != nil {
				output["private_ip_address"] = *props.PrivateIPAddress
			}

			if props.PublicIPAddress != nil && props.PublicIPAddress.Id != nil {
				output["public_ip_address_id"] = *props.PublicIPAddress.Id
			}

			if props.PrivateLinkConfiguration != nil && props.PrivateLinkConfiguration.Id != nil {
				configurationID, err := parse.ApplicationGatewayPrivateLinkConfigurationIDInsensitively(*props.PrivateLinkConfiguration.Id)
				if err != nil {
					return nil, err
				}
				output["private_link_configuration_name"] = configurationID.PrivateLinkConfigurationName
				output["private_link_configuration_id"] = *props.PrivateLinkConfiguration.Id
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayProbes(input *[]applicationgateways.ApplicationGatewayProbe) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			output["protocol"] = pointer.FromEnum(props.Protocol)

			if host := props.Host; host != nil {
				output["host"] = *host
			}

			if path := props.Path; path != nil {
				output["path"] = *path
			}

			if interval := props.Interval; interval != nil {
				output["interval"] = int(*interval)
			}

			if timeout := props.Timeout; timeout != nil {
				output["timeout"] = int(*timeout)
			}

			if threshold := props.UnhealthyThreshold; threshold != nil {
				output["unhealthy_threshold"] = int(*threshold)
			}

			port := 0
			if props.Port != nil {
				port = int(*props.Port)
			}
			output["port"] = port

			if pickHostNameFromBackendHTTPSettings := props.PickHostNameFromBackendHTTPSettings; pickHostNameFromBackendHTTPSettings != nil {
				output["pick_host_name_from_backend_http_settings"] = *pickHostNameFromBackendHTTPSettings
			}

			output["proxy_protocol_header_enabled"] = pointer.From(props.EnableProbeProxyProtocolHeader)

			if minServers := props.MinServers; minServers != nil {
				output["minimum_servers"] = int(*minServers)
			}

			matches := make([]any, 0)
			if match := props.Match; match != nil {
				matchConfig := map[string]any{}
				if body := match.Body; body != nil {
					matchConfig["body"] = *body
				}

				statusCodes := make([]any, 0)
				if match.StatusCodes != nil {
					for _, status := range *match.StatusCodes {
						statusCodes = append(statusCodes, status)
					}
				}
				matchConfig["status_code"] = statusCodes
				matches = append(matches, matchConfig)
			}
			output["match"] = matches
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayPrivateEndpoints(input *[]applicationgateways.ApplicationGatewayPrivateEndpointConnection) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]any{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}
		if endpoint.Id != nil {
			result["id"] = *endpoint.Id
		}
	}
	return results
}

func flattenApplicationGatewayPrivateLinkConfigurations(input *[]applicationgateways.ApplicationGatewayPrivateLinkConfiguration) []any {
	plConfigResults := make([]any, 0)
	if input == nil {
		return plConfigResults
	}

	for _, plConfig := range *input {
		plConfigResult := map[string]any{}
		if plConfig.Name != nil {
			plConfigResult["name"] = *plConfig.Name
		}
		if plConfig.Id != nil {
			plConfigResult["id"] = *plConfig.Id
		}
		ipConfigResults := make([]any, 0)
		if props := plConfig.Properties; props != nil {
			for _, ipConfig := range *props.IPConfigurations {
				ipConfigResult := map[string]any{}
				if ipConfig.Name != nil {
					ipConfigResult["name"] = *ipConfig.Name
				}
				if ipConfigProps := ipConfig.Properties; ipConfigProps != nil {
					if ipConfigProps.Subnet != nil {
						ipConfigResult["subnet_id"] = *ipConfigProps.Subnet.Id
					}
					if ipConfigProps.PrivateIPAddress != nil {
						ipConfigResult["private_ip_address"] = *ipConfigProps.PrivateIPAddress
					}
					ipConfigResult["private_ip_address_allocation"] = pointer.FromEnum(ipConfigProps.PrivateIPAllocationMethod)
					if ipConfigProps.Primary != nil {
						ipConfigResult["primary"] = *ipConfigProps.Primary
					}
					ipConfigResults = append(ipConfigResults, ipConfigResult)
				}
			}
		}
		plConfigResult["ip_configuration"] = ipConfigResults
		plConfigResults = append(plConfigResults, plConfigResult)
	}
	return plConfigResults
}

func flattenApplicationGatewayRequestRoutingRules(input *[]applicationgateways.ApplicationGatewayRequestRoutingRule) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]any{
				"rule_type": pointer.FromEnum(props.RuleType),
			}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if props.Priority != nil {
				output["priority"] = *props.Priority
			}

			if pool := props.BackendAddressPool; pool != nil {
				if pool.Id != nil {
					poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.Id)
					if err != nil {
						return nil, err
					}
					output["backend_address_pool_name"] = poolId.Name
					output["backend_address_pool_id"] = poolId.ID()
				}
			}

			if settings := props.BackendHTTPSettings; settings != nil {
				if settings.Id != nil {
					settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.Id)
					if err != nil {
						return nil, err
					}

					output["backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
					output["backend_http_settings_id"] = *settings.Id
				}
			}

			if listener := props.HTTPListener; listener != nil {
				if listener.Id != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.Id)
					if err != nil {
						return nil, err
					}
					output["http_listener_id"] = listenerId.ID()
					output["http_listener_name"] = listenerId.Name
				}
			}

			if pathMap := props.UrlPathMap; pathMap != nil {
				if pathMap.Id != nil {
					pathMapId, err := parse.UrlPathMapIDInsensitively(*pathMap.Id)
					if err != nil {
						return nil, err
					}
					output["url_path_map_name"] = pathMapId.Name
					output["url_path_map_id"] = pathMapId.ID()
				}
			}

			if redirect := props.RedirectConfiguration; redirect != nil {
				if redirect.Id != nil {
					redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
					if err != nil {
						return nil, err
					}
					output["redirect_configuration_name"] = redirectId.RedirectConfigurationName
					output["redirect_configuration_id"] = redirectId.ID()
				}
			}

			if rewrite := props.RewriteRuleSet; rewrite != nil {
				if rewrite.Id != nil {
					rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
					if err != nil {
						return nil, err
					}
					output["rewrite_rule_set_name"] = rewriteId.Name
					output["rewrite_rule_set_id"] = rewriteId.ID()
				}
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func flattenApplicationGatewayRoutingRules(input *[]applicationgateways.ApplicationGatewayRoutingRule) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]any{
				"id":       pointer.From(config.Id),
				"name":     pointer.From(config.Name),
				"priority": props.Priority,
			}

			if pool := props.BackendAddressPool; pool != nil && pool.Id != nil {
				poolId, err := parse.BackendAddressPoolIDInsensitively(pointer.From(pool.Id))
				if err != nil {
					return nil, err
				}
				output["backend_address_pool_name"] = poolId.Name
				output["backend_address_pool_id"] = poolId.ID()
			}

			if settings := props.BackendSettings; settings != nil && settings.Id != nil {
				settingsId, err := parse.BackendSettingsCollectionIDInsensitively(pointer.From(settings.Id))
				if err != nil {
					return nil, err
				}
				output["backend_name"] = settingsId.BackendSettingsCollectionName
				output["backend_id"] = settingsId.ID()
			}

			if listener := props.Listener; listener != nil && listener.Id != nil {
				listenerId, err := parse.ListenerIDInsensitively(pointer.From(listener.Id))
				if err != nil {
					return nil, err
				}
				output["listener_name"] = listenerId.Name
				output["listener_id"] = listenerId.ID()
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func flattenApplicationGatewayRewriteRuleSets(input *[]applicationgateways.ApplicationGatewayRewriteRuleSet) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]any{}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if rulesConfig := props.RewriteRules; rulesConfig != nil {
				rules := make([]any, 0)
				for _, rule := range *rulesConfig {
					ruleOutput := map[string]any{}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if rule.RuleSequence != nil {
						ruleOutput["rule_sequence"] = *rule.RuleSequence
					}

					conditions := make([]any, 0)
					if rule.Conditions != nil {
						for _, config := range *rule.Conditions {
							condition := map[string]any{}

							if config.Variable != nil {
								condition["variable"] = *config.Variable
							}

							if config.Pattern != nil {
								condition["pattern"] = *config.Pattern
							}

							if config.IgnoreCase != nil {
								condition["ignore_case"] = *config.IgnoreCase
							}

							if config.Negate != nil {
								condition["negate"] = *config.Negate
							}

							conditions = append(conditions, condition)
						}
					}
					ruleOutput["condition"] = conditions

					requestConfigs := make([]any, 0)
					responseConfigs := make([]any, 0)
					urlConfigs := make([]any, 0)

					if rule.ActionSet != nil {
						actionSet := *rule.ActionSet

						if actionSet.RequestHeaderConfigurations != nil {
							for _, config := range *actionSet.RequestHeaderConfigurations {
								requestConfig := map[string]any{}

								if config.HeaderName != nil {
									requestConfig["header_name"] = *config.HeaderName
								}

								if config.HeaderValue != nil {
									requestConfig["header_value"] = *config.HeaderValue
								}

								requestConfigs = append(requestConfigs, requestConfig)
							}
						}

						if actionSet.ResponseHeaderConfigurations != nil {
							for _, config := range *actionSet.ResponseHeaderConfigurations {
								responseConfig := map[string]any{}

								if config.HeaderName != nil {
									responseConfig["header_name"] = *config.HeaderName
								}

								if config.HeaderValue != nil {
									responseConfig["header_value"] = *config.HeaderValue
								}

								responseConfigs = append(responseConfigs, responseConfig)
							}
						}

						if actionSet.UrlConfiguration != nil {
							config := *actionSet.UrlConfiguration
							components := ""

							path := pointer.From(config.ModifiedPath)

							queryString := pointer.From(config.ModifiedQueryString)

							// `components` doesn't exist in the API - it appears to be purely a UI state in the Portal
							// as such we should consider removing this field in the future.
							if path == queryString {
								// used to represent `both`
								components = ""
							}
							if config.ModifiedQueryString != nil && config.ModifiedPath == nil {
								components = "query_string_only"
							}
							if config.ModifiedQueryString == nil && config.ModifiedPath != nil {
								components = "path_only"
							}

							urlConfigs = append(urlConfigs, map[string]any{
								"components":   components,
								"query_string": queryString,
								"path":         path,
								"reroute":      pointer.From(config.Reroute),
							})
						}
					}
					ruleOutput["request_header_configuration"] = requestConfigs
					ruleOutput["response_header_configuration"] = responseConfigs
					ruleOutput["url"] = urlConfigs

					rules = append(rules, ruleOutput)
				}
				output["rewrite_rule"] = rules
			}
			results = append(results, output)
		}
	}

	return results
}

func flattenApplicationGatewayRedirectConfigurations(input *[]applicationgateways.ApplicationGatewayRedirectConfiguration) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]any{
				"redirect_type": pointer.FromEnum(props.RedirectType),
			}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if listener := props.TargetListener; listener != nil {
				if listener.Id != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.Id)
					if err != nil {
						return nil, err
					}
					output["target_listener_name"] = listenerId.Name
					output["target_listener_id"] = listenerId.ID()
				}
			}

			if props.TargetURL != nil {
				output["target_url"] = *props.TargetURL
			}

			if props.IncludePath != nil {
				output["include_path"] = *props.IncludePath
			}

			if props.IncludeQueryString != nil {
				output["include_query_string"] = *props.IncludeQueryString
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func flattenApplicationGatewayAutoscaleConfiguration(input *applicationgateways.ApplicationGatewayAutoscaleConfiguration) []any {
	result := make(map[string]any)
	if input == nil {
		return []any{}
	}

	result["min_capacity"] = input.MinCapacity
	if input.MaxCapacity != nil {
		result["max_capacity"] = *input.MaxCapacity
	}

	return []any{result}
}

func flattenApplicationGatewaySku(input *applicationgateways.ApplicationGatewaySku) []any {
	result := make(map[string]any)

	result["name"] = pointer.FromEnum(input.Name)
	result["tier"] = pointer.FromEnum(input.Tier)
	if input.Capacity != nil {
		result["capacity"] = int(*input.Capacity)
	}

	return []any{result}
}

func flattenApplicationGatewaySslCertificates(input *[]applicationgateways.ApplicationGatewaySslCertificate, d *pluginsdk.ResourceData) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}
		if v.Name == nil {
			continue
		}

		name := *v.Name

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = name

		if props := v.Properties; props != nil {
			if data := props.PublicCertData; data != nil {
				output["public_cert_data"] = *data
			}

			if kvsid := props.KeyVaultSecretId; kvsid != nil {
				output["key_vault_secret_id"] = *kvsid
			}
		}

		// since the certificate data isn't returned we have to load it from the same index
		if existing, ok := d.GetOk("ssl_certificate"); ok && existing != nil {
			existingVals := existing.(*schema.Set).List()

			for _, existingVal := range existingVals {
				existingCerts := existingVal.(map[string]any)
				existingName := existingCerts["name"].(string)

				if name == existingName {
					if data := existingCerts["data"]; data != nil {
						output["data"] = base64.EncodeIfNot(data.(string))
					}

					if password := existingCerts["password"]; password != nil {
						output["password"] = password.(string)
					}
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayTrustedClientCertificates(input *[]applicationgateways.ApplicationGatewayTrustedClientCertificate) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}
		if v.Name == nil {
			continue
		}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = *v.Name

		if props := v.Properties; props != nil {
			if data := props.Data; data != nil {
				output["data"] = *data
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewaySslProfiles(input *[]applicationgateways.ApplicationGatewaySslProfile) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{}
		if v.Name == nil {
			continue
		}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = *v.Name
		output["ssl_policy"] = flattenApplicationGatewaySslPolicy(v.Properties.SslPolicy)

		clientAuthenticationMode := string(applicationgateways.ApplicationGatewayClientAuthVerificationModesStrict)
		verifyClientCertIssuerDn := false
		verifyClientCertificateRevocation := ""

		if props := v.Properties; props != nil {
			if props.ClientAuthConfiguration != nil {
				verifyClientCertIssuerDn = pointer.From(props.ClientAuthConfiguration.VerifyClientCertIssuerDN)
				if *props.ClientAuthConfiguration.VerifyClientRevocation != applicationgateways.ApplicationGatewayClientRevocationOptionsNone {
					verifyClientCertificateRevocation = pointer.FromEnum(props.ClientAuthConfiguration.VerifyClientRevocation)
				}
				if props.ClientAuthConfiguration.VerifyClientAuthMode != nil {
					clientAuthenticationMode = pointer.FromEnum(props.ClientAuthConfiguration.VerifyClientAuthMode)
				}
			}

			trustedClientCertificateNames := make([]any, 0)
			if certs := props.TrustedClientCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.TrustedClientCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					trustedClientCertificateNames = append(trustedClientCertificateNames, certId.Name)
				}
			}
			output["trusted_client_certificate_names"] = trustedClientCertificateNames
			output["client_authentication_mode"] = clientAuthenticationMode
			output["verify_client_certificate_issuer_dn"] = verifyClientCertIssuerDn
			output["verify_client_certificate_revocation"] = verifyClientCertificateRevocation
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayURLPathMaps(input *[]applicationgateways.ApplicationGatewayURLPathMap) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]any{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if backendPool := props.DefaultBackendAddressPool; backendPool != nil && backendPool.Id != nil {
				poolId, err := parse.BackendAddressPoolIDInsensitively(*backendPool.Id)
				if err != nil {
					return nil, err
				}
				output["default_backend_address_pool_name"] = poolId.Name
				output["default_backend_address_pool_id"] = poolId.ID()
			}

			if settings := props.DefaultBackendHTTPSettings; settings != nil && settings.Id != nil {
				settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.Id)
				if err != nil {
					return nil, err
				}
				output["default_backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
				output["default_backend_http_settings_id"] = settingsId.ID()
			}

			if redirect := props.DefaultRedirectConfiguration; redirect != nil && redirect.Id != nil {
				redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
				if err != nil {
					return nil, err
				}
				output["default_redirect_configuration_name"] = redirectId.RedirectConfigurationName
				output["default_redirect_configuration_id"] = redirectId.ID()
			}

			if rewrite := props.DefaultRewriteRuleSet; rewrite != nil && rewrite.Id != nil {
				rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
				if err != nil {
					return nil, err
				}
				output["default_rewrite_rule_set_name"] = rewriteId.Name
				output["default_rewrite_rule_set_id"] = rewriteId.ID()
			}

			pathRules := make([]any, 0)
			if rules := props.PathRules; rules != nil {
				for _, rule := range *rules {
					ruleOutput := map[string]any{}

					if rule.Id != nil {
						ruleOutput["id"] = *rule.Id
					}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if ruleProps := rule.Properties; ruleProps != nil {
						if pool := ruleProps.BackendAddressPool; pool != nil && pool.Id != nil {
							poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_address_pool_name"] = poolId.Name
							ruleOutput["backend_address_pool_id"] = poolId.ID()
						}

						if backend := ruleProps.BackendHTTPSettings; backend != nil && backend.Id != nil {
							backendId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*backend.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_http_settings_name"] = backendId.BackendHttpSettingsCollectionName
							ruleOutput["backend_http_settings_id"] = backendId.ID()
						}

						if redirect := ruleProps.RedirectConfiguration; redirect != nil && redirect.Id != nil {
							redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["redirect_configuration_name"] = redirectId.RedirectConfigurationName
							ruleOutput["redirect_configuration_id"] = redirectId.ID()
						}

						if rewrite := ruleProps.RewriteRuleSet; rewrite != nil && rewrite.Id != nil {
							rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["rewrite_rule_set_name"] = rewriteId.Name
							ruleOutput["rewrite_rule_set_id"] = rewriteId.ID()
						}

						if fwp := ruleProps.FirewallPolicy; fwp != nil && fwp.Id != nil {
							policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["firewall_policy_id"] = policyId.ID()
						}

						pathOutputs := make([]any, 0)
						if paths := ruleProps.Paths; paths != nil {
							for _, rulePath := range *paths {
								pathOutputs = append(pathOutputs, rulePath)
							}
						}
						ruleOutput["paths"] = pathOutputs
					}

					pathRules = append(pathRules, ruleOutput)
				}
			}
			output["path_rule"] = pathRules
		}

		results = append(results, output)
	}

	return results, nil
}

func flattenApplicationGatewayWafConfig(input *applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	output := make(map[string]any)

	output["enabled"] = input.Enabled
	output["firewall_mode"] = string(input.FirewallMode)
	output["rule_set_type"] = input.RuleSetType
	output["rule_set_version"] = input.RuleSetVersion

	if input.DisabledRuleGroups != nil {
		output["disabled_rule_group"] = flattenApplicationGateWayDisabledRuleGroups(input.DisabledRuleGroups)
	}

	if input.FileUploadLimitInMb != nil {
		output["file_upload_limit_mb"] = int(*input.FileUploadLimitInMb)
	}

	if input.RequestBodyCheck != nil {
		output["request_body_check"] = *input.RequestBodyCheck
	}

	if input.MaxRequestBodySizeInKb != nil {
		output["max_request_body_size_kb"] = int(*input.MaxRequestBodySizeInKb)
	}

	if input.Exclusions != nil {
		output["exclusion"] = flattenApplicationGatewayFirewallExclusion(input.Exclusions)
	}
	results = append(results, output)

	return results
}

func flattenApplicationGateWayDisabledRuleGroups(input *[]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup) []any {
	ruleGroups := make([]any, 0)
	for _, ruleGroup := range *input {
		ruleGroupOutput := map[string]any{}

		ruleGroupOutput["rule_group_name"] = ruleGroup.RuleGroupName

		ruleOutputs := make([]any, 0)
		if rules := ruleGroup.Rules; rules != nil {
			for _, rule := range *rules {
				ruleOutputs = append(ruleOutputs, rule)
			}
		}
		ruleGroupOutput["rules"] = ruleOutputs

		ruleGroups = append(ruleGroups, ruleGroupOutput)
	}
	return ruleGroups
}

func flattenApplicationGatewayFirewallExclusion(input *[]applicationgateways.ApplicationGatewayFirewallExclusion) []any {
	exclusionLists := make([]any, 0)
	for _, exclusionList := range *input {
		exclusionListOutput := map[string]any{}
		exclusionListOutput["match_variable"] = exclusionList.MatchVariable
		exclusionListOutput["selector_match_operator"] = exclusionList.SelectorMatchOperator
		exclusionListOutput["selector"] = exclusionList.Selector
		exclusionLists = append(exclusionLists, exclusionListOutput)
	}
	return exclusionLists
}

func flattenApplicationGatewayCustomErrorConfigurations(input *[]applicationgateways.ApplicationGatewayCustomError) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]any{}

		output["status_code"] = v.StatusCode

		if v.CustomErrorPageURL != nil {
			output["custom_error_page_url"] = *v.CustomErrorPageURL
		}

		results = append(results, output)
	}

	return results
}

// backendHTTPSettingsUnderUsedNames keeps each backend HTTP setting under the cookie affinity name it already uses in state, as neither name can be
// Computed inside a set. Settings not in state yet use the new name. Only needed until 6.0 removes `cookie_based_affinity`.
func backendHTTPSettingsUnderUsedNames(d *pluginsdk.ResourceData, settings []any) []any {
	usesOldName := make(map[string]bool)
	if existing, ok := d.Get("backend_http_settings").(*pluginsdk.Set); ok {
		for _, raw := range existing.List() {
			if m, ok := raw.(map[string]any); ok && m["cookie_based_affinity"] != "" {
				usesOldName[m["name"].(string)] = true
			}
		}
	}

	for _, raw := range settings {
		m := raw.(map[string]any)
		if usesOldName[m["name"].(string)] {
			m["cookie_based_affinity_enabled"] = false
		} else {
			m["cookie_based_affinity"] = ""
		}
	}

	return settings
}
