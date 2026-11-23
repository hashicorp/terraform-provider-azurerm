// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApplicationGatewayTrustedRootCertificates(certs []any) (*[]applicationgateways.ApplicationGatewayTrustedRootCertificate, error) {
	results := make([]applicationgateways.ApplicationGatewayTrustedRootCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		data := v["data"].(string)
		kvsid := v["key_vault_secret_id"].(string)

		output := applicationgateways.ApplicationGatewayTrustedRootCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayTrustedRootCertificatePropertiesFormat{},
		}

		switch {
		case data != "" && kvsid != "":
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		case data != "":
			output.Properties.Data = pointer.To(base64.EncodeIfNot(data))
		case kvsid != "":
			output.Properties.KeyVaultSecretId = pointer.To(kvsid)
		default:
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandApplicationGatewayBackendAddressPools(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayBackendAddressPool {
	vs := d.Get("backend_address_pool").(*schema.Set).List()
	results := make([]applicationgateways.ApplicationGatewayBackendAddressPool, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)
		backendAddresses := make([]applicationgateways.ApplicationGatewayBackendAddress, 0)

		if fqdnsConfig, ok := v["fqdns"]; ok {
			fqdns := fqdnsConfig.(*schema.Set).List()
			for _, ip := range fqdns {
				backendAddresses = append(backendAddresses, applicationgateways.ApplicationGatewayBackendAddress{
					Fqdn: pointer.To(ip.(string)),
				})
			}
		}

		if ipAddressesConfig, ok := v["ip_addresses"]; ok {
			ipAddresses := ipAddressesConfig.(*schema.Set).List()

			for _, ip := range ipAddresses {
				backendAddresses = append(backendAddresses, applicationgateways.ApplicationGatewayBackendAddress{
					IPAddress: pointer.To(ip.(string)),
				})
			}
		}

		name := v["name"].(string)
		output := applicationgateways.ApplicationGatewayBackendAddressPool{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayBackendAddressPoolPropertiesFormat{
				BackendAddresses: &backendAddresses,
			},
		}

		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayBackendHTTPSettings(input []any, gatewayID string) *[]applicationgateways.ApplicationGatewayBackendHTTPSettings {
	results := make([]applicationgateways.ApplicationGatewayBackendHTTPSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		path := v["path"].(string)
		port := int64(v["port"].(int))
		protocol := v["protocol"].(string)
		cookieBasedAffinity := v["cookie_based_affinity"].(string)
		pickHostNameFromBackendAddress := v["pick_host_name_from_backend_address"].(bool)
		requestTimeout := int64(v["request_timeout"].(int))

		setting := applicationgateways.ApplicationGatewayBackendHTTPSettings{
			Name: pointer.To(v["name"].(string)),
			Properties: &applicationgateways.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
				ConnectionDraining:             expandApplicationGatewayConnectionDraining(v),
				CookieBasedAffinity:            pointer.ToEnum[applicationgateways.ApplicationGatewayCookieBasedAffinity](cookieBasedAffinity),
				DedicatedBackendConnection:     pointer.To(v["dedicated_backend_connection_enabled"].(bool)),
				Path:                           pointer.To(path),
				PickHostNameFromBackendAddress: pointer.To(pickHostNameFromBackendAddress),
				Port:                           pointer.To(port),
				Protocol:                       pointer.ToEnum[applicationgateways.ApplicationGatewayProtocol](protocol),
				RequestTimeout:                 pointer.To(requestTimeout),
				SniName:                        pointer.To(v["sni_name"].(string)),
				ValidateCertChainAndExpiry:     pointer.To(v["certificate_chain_validation_enabled"].(bool)),
				ValidateSNI:                    pointer.To(v["sni_validation_enabled"].(bool)),
			},
		}

		hostName := v["host_name"].(string)
		if hostName != "" {
			setting.Properties.HostName = pointer.To(hostName)
		}

		affinityCookieName := v["affinity_cookie_name"].(string)
		if affinityCookieName != "" {
			setting.Properties.AffinityCookieName = pointer.To(affinityCookieName)
		}

		if v["trusted_root_certificate_names"] != nil {
			trustedRootCertNames := v["trusted_root_certificate_names"].([]any)
			trustedRootCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawTrustedRootCertName := range trustedRootCertNames {
				trustedRootCertName := rawTrustedRootCertName.(string)
				trustedRootCertID := fmt.Sprintf("%s/trustedRootCertificates/%s", gatewayID, trustedRootCertName)
				trustedRootCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(trustedRootCertID),
				}

				trustedRootCertSubResources = append(trustedRootCertSubResources, trustedRootCertSubResource)
			}

			setting.Properties.TrustedRootCertificates = &trustedRootCertSubResources
		}

		probeName := v["probe_name"].(string)
		if probeName != "" {
			probeID := fmt.Sprintf("%s/probes/%s", gatewayID, probeName)
			setting.Properties.Probe = &applicationgateways.SubResource{
				Id: pointer.To(probeID),
			}
		}

		results = append(results, setting)
	}

	return &results
}

func expandApplicationGatewayConnectionDraining(d map[string]any) *applicationgateways.ApplicationGatewayConnectionDraining {
	connectionsRaw := d["connection_draining"].([]any)

	if len(connectionsRaw) == 0 {
		return nil
	}

	connectionRaw := connectionsRaw[0].(map[string]any)

	return &applicationgateways.ApplicationGatewayConnectionDraining{
		Enabled:           connectionRaw["enabled"].(bool),
		DrainTimeoutInSec: int64(connectionRaw["drain_timeout_sec"].(int)),
	}
}

func expandApplicationGatewayBackendSettings(input []any, appGwID applicationgateways.ApplicationGatewayId) *[]applicationgateways.ApplicationGatewayBackendSettings {
	results := make([]applicationgateways.ApplicationGatewayBackendSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		setting := applicationgateways.ApplicationGatewayBackendSettings{
			Name: pointer.To(v["name"].(string)),
			Properties: &applicationgateways.ApplicationGatewayBackendSettingsPropertiesFormat{
				Port:                         pointer.To(int64(v["port"].(int))),
				Protocol:                     pointer.ToEnum[applicationgateways.ApplicationGatewayProtocol](v["protocol"].(string)),
				Timeout:                      pointer.To(int64(v["timeout_in_seconds"].(int))),
				EnableL4ClientIPPreservation: pointer.To(v["client_ip_preservation_enabled"].(bool)),
			},
		}

		if hostName := v["host_name"].(string); hostName != "" {
			setting.Properties.HostName = pointer.To(hostName)
		}

		if trustedRootCertNames := v["trusted_root_certificate_names"].([]any); len(trustedRootCertNames) != 0 {
			trustedRootCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawTrustedRootCertName := range trustedRootCertNames {
				trustedRootCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(parse.NewTrustedRootCertificateID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, rawTrustedRootCertName.(string)).ID()),
				}

				trustedRootCertSubResources = append(trustedRootCertSubResources, trustedRootCertSubResource)
			}

			setting.Properties.TrustedRootCertificates = &trustedRootCertSubResources
		}

		if probeName := v["probe_name"].(string); probeName != "" {
			setting.Properties.Probe = &applicationgateways.SubResource{
				Id: pointer.To(parse.NewProbeID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, probeName).ID()),
			}
		}

		results = append(results, setting)
	}

	return &results
}

func expandApplicationGatewaySslPolicy(vs []any) *applicationgateways.ApplicationGatewaySslPolicy {
	policy := applicationgateways.ApplicationGatewaySslPolicy{}
	disabledSSLProtocols := make([]applicationgateways.ApplicationGatewaySslProtocol, 0)

	if len(vs) > 0 && vs[0] != nil {
		v := vs[0].(map[string]any)
		policyType := applicationgateways.ApplicationGatewaySslPolicyType(v["policy_type"].(string))

		for _, policy := range v["disabled_protocols"].([]any) {
			disabledSSLProtocols = append(disabledSSLProtocols, applicationgateways.ApplicationGatewaySslProtocol(policy.(string)))
		}

		switch policyType {
		case applicationgateways.ApplicationGatewaySslPolicyTypePredefined:
			policyName := applicationgateways.ApplicationGatewaySslPolicyName(v["policy_name"].(string))
			policy = applicationgateways.ApplicationGatewaySslPolicy{
				PolicyType: pointer.To(policyType),
				PolicyName: pointer.To(policyName),
			}
		case applicationgateways.ApplicationGatewaySslPolicyTypeCustom, applicationgateways.ApplicationGatewaySslPolicyTypeCustomVTwo:
			minProtocolVersion := applicationgateways.ApplicationGatewaySslProtocol(v["min_protocol_version"].(string))
			cipherSuites := make([]applicationgateways.ApplicationGatewaySslCipherSuite, 0)

			for _, cipherSuite := range v["cipher_suites"].([]any) {
				cipherSuites = append(cipherSuites, applicationgateways.ApplicationGatewaySslCipherSuite(cipherSuite.(string)))
			}

			policy = applicationgateways.ApplicationGatewaySslPolicy{
				PolicyType:         pointer.To(policyType),
				MinProtocolVersion: pointer.To(minProtocolVersion),
				CipherSuites:       &cipherSuites,
			}
		}
	}

	if len(disabledSSLProtocols) > 0 {
		policy = applicationgateways.ApplicationGatewaySslPolicy{
			DisabledSslProtocols: &disabledSSLProtocols,
		}
	}

	return &policy
}

func expandApplicationGatewayHTTPListeners(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayHTTPListener, error) {
	vs := d.Get("http_listener").(*schema.Set).List()

	results := make([]applicationgateways.ApplicationGatewayHTTPListener, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		frontendIPConfigName := v["frontend_ip_configuration_name"].(string)
		frontendPortName := v["frontend_port_name"].(string)
		protocol := v["protocol"].(string)
		requireSNI := v["require_sni"].(bool)
		sslProfileName := v["ssl_profile_name"].(string)

		frontendIPConfigID := fmt.Sprintf("%s/frontendIPConfigurations/%s", gatewayID, frontendIPConfigName)
		frontendPortID := fmt.Sprintf("%s/frontendPorts/%s", gatewayID, frontendPortName)
		firewallPolicyID := v["firewall_policy_id"].(string)

		customErrorConfigurations := expandApplicationGatewayCustomErrorConfigurations(v["custom_error_configuration"].([]any))

		listener := applicationgateways.ApplicationGatewayHTTPListener{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayHTTPListenerPropertiesFormat{
				FrontendIPConfiguration: &applicationgateways.SubResource{
					Id: pointer.To(frontendIPConfigID),
				},
				FrontendPort: &applicationgateways.SubResource{
					Id: pointer.To(frontendPortID),
				},
				Protocol:                    pointer.ToEnum[applicationgateways.ApplicationGatewayProtocol](protocol),
				RequireServerNameIndication: pointer.To(requireSNI),
				CustomErrorConfigurations:   customErrorConfigurations,
			},
		}

		host := v["host_name"].(string)
		hosts := v["host_names"].(*pluginsdk.Set).List()

		if host != "" && len(hosts) > 0 {
			return nil, fmt.Errorf("`host_name` and `host_names` cannot be specified together")
		}

		if host != "" {
			listener.Properties.HostName = &host
		}

		if len(hosts) > 0 {
			listener.Properties.HostNames = pluginsdk.ExpandStringSlice(hosts)
		}

		if sslCertName := v["ssl_certificate_name"].(string); sslCertName != "" {
			certID := fmt.Sprintf("%s/sslCertificates/%s", gatewayID, sslCertName)
			listener.Properties.SslCertificate = &applicationgateways.SubResource{
				Id: pointer.To(certID),
			}
		}

		if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
			listener.Properties.FirewallPolicy = &applicationgateways.SubResource{
				Id: pointer.To(firewallPolicyID),
			}
		}

		if sslProfileName != "" && len(sslProfileName) > 0 {
			sslProfileID := fmt.Sprintf("%s/sslProfiles/%s", gatewayID, sslProfileName)
			listener.Properties.SslProfile = &applicationgateways.SubResource{
				Id: pointer.To(sslProfileID),
			}
		}

		results = append(results, listener)
	}

	return &results, nil
}

func expandApplicationGatewayListeners(input []any, appGwID applicationgateways.ApplicationGatewayId) *[]applicationgateways.ApplicationGatewayListener {
	results := make([]applicationgateways.ApplicationGatewayListener, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		listener := applicationgateways.ApplicationGatewayListener{
			Name: pointer.To(v["name"].(string)),
			Properties: &applicationgateways.ApplicationGatewayListenerPropertiesFormat{
				FrontendIPConfiguration: &applicationgateways.SubResource{
					Id: pointer.To(parse.NewFrontendIPConfigurationID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, v["frontend_ip_configuration_name"].(string)).ID()),
				},
				FrontendPort: &applicationgateways.SubResource{
					Id: pointer.To(parse.NewFrontendPortID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, v["frontend_port_name"].(string)).ID()),
				},
				Protocol: pointer.ToEnum[applicationgateways.ApplicationGatewayProtocol](v["protocol"].(string)),
			},
		}

		if hosts := v["host_names"].(*pluginsdk.Set).List(); len(hosts) > 0 {
			listener.Properties.HostNames = pluginsdk.ExpandStringSlice(hosts)
		}

		if sslCertName := v["ssl_certificate_name"].(string); sslCertName != "" {
			listener.Properties.SslCertificate = &applicationgateways.SubResource{
				Id: pointer.To(parse.NewSslCertificateID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, sslCertName).ID()),
			}
		}

		if sslProfileName := v["ssl_profile_name"].(string); sslProfileName != "" {
			listener.Properties.SslProfile = &applicationgateways.SubResource{
				Id: pointer.To(parse.NewSslProfileID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, sslProfileName).ID()),
			}
		}

		results = append(results, listener)
	}

	return &results
}

func expandApplicationGatewayIPConfigurations(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayIPConfiguration, bool) {
	vs := d.Get("gateway_ip_configuration").([]any)
	results := make([]applicationgateways.ApplicationGatewayIPConfiguration, 0)
	stopApplicationGateway := false

	for _, configRaw := range vs {
		data := configRaw.(map[string]any)

		name := data["name"].(string)
		subnetID := data["subnet_id"].(string)

		output := applicationgateways.ApplicationGatewayIPConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayIPConfigurationPropertiesFormat{
				Subnet: &applicationgateways.SubResource{
					Id: pointer.To(subnetID),
				},
			},
		}
		results = append(results, output)
	}

	if d.HasChange("gateway_ip_configuration") {
		oldRaw, newRaw := d.GetChange("gateway_ip_configuration")
		oldVS := oldRaw.([]any)
		newVS := newRaw.([]any)

		// If we're creating the application gateway return the current gateway ip configuration.
		if len(oldVS) == 0 {
			return &results, false
		}

		// The application gateway needs to be stopped if a gateway ip configuration is added or removed
		if len(oldVS) != len(newVS) {
			return &results, true
		}

		for i, configRaw := range newVS {
			newData := configRaw.(map[string]any)
			oldData := oldVS[i].(map[string]any)

			newSubnetID := newData["subnet_id"].(string)
			oldSubnetID := oldData["subnet_id"].(string)
			// The application gateway needs to be shutdown if the subnet ids don't match
			if newSubnetID != oldSubnetID {
				stopApplicationGateway = true
			}
		}
	}

	return &results, stopApplicationGateway
}

func expandApplicationGatewayGlobalConfiguration(input []any) *applicationgateways.ApplicationGatewayGlobalConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	return &applicationgateways.ApplicationGatewayGlobalConfiguration{
		EnableRequestBuffering:  pointer.To(v["request_buffering_enabled"].(bool)),
		EnableResponseBuffering: pointer.To(v["response_buffering_enabled"].(bool)),
	}
}

func expandApplicationGatewayFrontendPorts(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayFrontendPort {
	vs := d.Get("frontend_port").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayFrontendPort, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		port := int64(v["port"].(int))

		output := applicationgateways.ApplicationGatewayFrontendPort{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayFrontendPortPropertiesFormat{
				Port: pointer.To(port),
			},
		}
		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayFrontendIPConfigurations(d *pluginsdk.ResourceData, gatewayID string) *[]applicationgateways.ApplicationGatewayFrontendIPConfiguration {
	vs := d.Get("frontend_ip_configuration").([]any)
	results := make([]applicationgateways.ApplicationGatewayFrontendIPConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		properties := applicationgateways.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{}

		if val := v["subnet_id"].(string); val != "" {
			properties.Subnet = &applicationgateways.SubResource{
				Id: pointer.To(val),
			}
		}

		if val := v["private_ip_address_allocation"].(string); val != "" {
			properties.PrivateIPAllocationMethod = pointer.ToEnum[applicationgateways.IPAllocationMethod](val)
		}

		if val := v["private_ip_address"].(string); val != "" {
			properties.PrivateIPAddress = pointer.To(val)
		}

		if val := v["public_ip_address_id"].(string); val != "" {
			properties.PublicIPAddress = &applicationgateways.SubResource{
				Id: pointer.To(val),
			}
		}

		if val := v["private_link_configuration_name"].(string); val != "" {
			privateLinkConfigurationID := fmt.Sprintf("%s/privateLinkConfigurations/%s", gatewayID, val)
			properties.PrivateLinkConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(privateLinkConfigurationID),
			}
		}

		name := v["name"].(string)
		output := applicationgateways.ApplicationGatewayFrontendIPConfiguration{
			Name:       pointer.To(name),
			Properties: &properties,
		}

		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayProbes(input []any) *[]applicationgateways.ApplicationGatewayProbe {
	results := make([]applicationgateways.ApplicationGatewayProbe, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		host := v["host"].(string)
		interval := int64(v["interval"].(int))
		minServers := int64(v["minimum_servers"].(int))
		name := v["name"].(string)
		probePath := v["path"].(string)
		protocol := v["protocol"].(string)
		port := int64(v["port"].(int))
		timeout := int64(v["timeout"].(int))
		unhealthyThreshold := int64(v["unhealthy_threshold"].(int))
		pickHostNameFromBackendHTTPSettings := v["pick_host_name_from_backend_http_settings"].(bool)

		output := applicationgateways.ApplicationGatewayProbe{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayProbePropertiesFormat{
				Host:                                pointer.To(host),
				Interval:                            pointer.To(interval),
				MinServers:                          pointer.To(minServers),
				Path:                                pointer.To(probePath),
				Protocol:                            pointer.ToEnum[applicationgateways.ApplicationGatewayProtocol](protocol),
				Timeout:                             pointer.To(timeout),
				UnhealthyThreshold:                  pointer.To(unhealthyThreshold),
				PickHostNameFromBackendHTTPSettings: pointer.To(pickHostNameFromBackendHTTPSettings),
				EnableProbeProxyProtocolHeader:      pointer.To(v["proxy_protocol_header_enabled"].(bool)),
			},
		}

		matchConfigs := v["match"].([]any)
		if len(matchConfigs) > 0 {
			matchBody := ""
			outputMatch := &applicationgateways.ApplicationGatewayProbeHealthResponseMatch{}
			if matchConfigs[0] != nil {
				match := matchConfigs[0].(map[string]any)
				matchBody = match["body"].(string)

				statusCodes := make([]string, 0)
				for _, statusCode := range match["status_code"].([]any) {
					statusCodes = append(statusCodes, statusCode.(string))
				}
				outputMatch.StatusCodes = &statusCodes
			}
			outputMatch.Body = pointer.To(matchBody)
			output.Properties.Match = outputMatch
		}

		if port != 0 {
			output.Properties.Port = pointer.To(port)
		}

		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayPrivateLinkConfigurations(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayPrivateLinkConfiguration {
	vs := d.Get("private_link_configuration").(*pluginsdk.Set).List()
	plConfigResults := make([]applicationgateways.ApplicationGatewayPrivateLinkConfiguration, 0)

	for _, rawPl := range vs {
		v := rawPl.(map[string]any)
		name := v["name"].(string)
		ipConfigurations := v["ip_configuration"].([]any)
		ipConfigurationResults := make([]applicationgateways.ApplicationGatewayPrivateLinkIPConfiguration, 0)
		for _, rawIp := range ipConfigurations {
			v := rawIp.(map[string]any)
			name := v["name"].(string)
			subnetId := v["subnet_id"].(string)
			ipConfiguration := applicationgateways.ApplicationGatewayPrivateLinkIPConfiguration{
				Name: pointer.To(name),
				Properties: &applicationgateways.ApplicationGatewayPrivateLinkIPConfigurationProperties{
					Primary: pointer.To(v["primary"].(bool)),
					Subnet: &applicationgateways.SubResource{
						Id: pointer.To(subnetId),
					},
				},
			}
			if privateIpAddress := v["private_ip_address"].(string); privateIpAddress != "" {
				ipConfiguration.Properties.PrivateIPAddress = pointer.To(privateIpAddress)
			}
			if privateIpAddressAllocation := v["private_ip_address_allocation"].(string); privateIpAddressAllocation != "" {
				ipConfiguration.Properties.PrivateIPAllocationMethod = pointer.ToEnum[applicationgateways.IPAllocationMethod](privateIpAddressAllocation)
			}
			ipConfigurationResults = append(ipConfigurationResults, ipConfiguration)
		}

		configuration := applicationgateways.ApplicationGatewayPrivateLinkConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayPrivateLinkConfigurationProperties{
				IPConfigurations: &ipConfigurationResults,
			},
		}
		plConfigResults = append(plConfigResults, configuration)
	}

	return &plConfigResults
}

func expandApplicationGatewayRequestRoutingRules(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayRequestRoutingRule, error) {
	vs := d.Get("request_routing_rule").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayRequestRoutingRule, 0)
	priorityset := false

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		ruleType := v["rule_type"].(string)
		httpListenerName := v["http_listener_name"].(string)
		httpListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, httpListenerName)
		backendAddressPoolName := v["backend_address_pool_name"].(string)
		backendHTTPSettingsName := v["backend_http_settings_name"].(string)
		redirectConfigName := v["redirect_configuration_name"].(string)
		priority := int64(v["priority"].(int))

		rule := applicationgateways.ApplicationGatewayRequestRoutingRule{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRequestRoutingRulePropertiesFormat{
				RuleType: pointer.ToEnum[applicationgateways.ApplicationGatewayRequestRoutingRuleType](ruleType),
				HTTPListener: &applicationgateways.SubResource{
					Id: pointer.To(httpListenerID),
				},
			},
		}

		if backendAddressPoolName != "" && redirectConfigName != "" {
			return nil, errors.New("conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
		}

		if backendHTTPSettingsName != "" && redirectConfigName != "" {
			return nil, errors.New("conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
		}

		if backendAddressPoolName != "" {
			backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
			rule.Properties.BackendAddressPool = &applicationgateways.SubResource{
				Id: pointer.To(backendAddressPoolID),
			}
		}

		if backendHTTPSettingsName != "" {
			backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
			rule.Properties.BackendHTTPSettings = &applicationgateways.SubResource{
				Id: pointer.To(backendHTTPSettingsID),
			}
		}

		if redirectConfigName != "" {
			redirectConfigID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigName)
			rule.Properties.RedirectConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(redirectConfigID),
			}
		}

		if urlPathMapName := v["url_path_map_name"].(string); urlPathMapName != "" {
			urlPathMapID := fmt.Sprintf("%s/urlPathMaps/%s", gatewayID, urlPathMapName)
			rule.Properties.UrlPathMap = &applicationgateways.SubResource{
				Id: pointer.To(urlPathMapID),
			}
		}

		if rewriteRuleSetName := v["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
			rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
			rule.Properties.RewriteRuleSet = &applicationgateways.SubResource{
				Id: pointer.To(rewriteRuleSetID),
			}
		}

		if priority != 0 {
			rule.Properties.Priority = &priority
			priorityset = true
		}

		results = append(results, rule)
	}

	if priorityset {
		for _, rule := range results {
			if rule.Properties.Priority == nil {
				return nil, errors.New("if you wish to use rule priority, you will have to specify rule-priority field values for all the existing request routing rules")
			}
		}
	}

	return &results, nil
}

func expandApplicationGatewayRoutingRules(input []any, appGwID applicationgateways.ApplicationGatewayId) *[]applicationgateways.ApplicationGatewayRoutingRule {
	results := make([]applicationgateways.ApplicationGatewayRoutingRule, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		rule := applicationgateways.ApplicationGatewayRoutingRule{
			Name: pointer.To(v["name"].(string)),
			Properties: &applicationgateways.ApplicationGatewayRoutingRulePropertiesFormat{
				RuleType: pointer.To(applicationgateways.ApplicationGatewayRequestRoutingRuleTypeBasic),
				Listener: &applicationgateways.SubResource{
					Id: pointer.To(parse.NewListenerID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, v["listener_name"].(string)).ID()),
				},
				BackendAddressPool: &applicationgateways.SubResource{
					Id: pointer.To(parse.NewBackendAddressPoolID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, v["backend_address_pool_name"].(string)).ID()),
				},
				BackendSettings: &applicationgateways.SubResource{
					Id: pointer.To(parse.NewBackendSettingsCollectionID(appGwID.SubscriptionId, appGwID.ResourceGroupName, appGwID.ApplicationGatewayName, v["backend_name"].(string)).ID()),
				},
				Priority: int64(v["priority"].(int)),
			},
		}

		results = append(results, rule)
	}

	return &results
}

func expandApplicationGatewayRewriteRuleSets(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayRewriteRuleSet, error) {
	vs := d.Get("rewrite_rule_set").([]any)
	ruleSets := make([]applicationgateways.ApplicationGatewayRewriteRuleSet, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)
		rules := make([]applicationgateways.ApplicationGatewayRewriteRule, 0)

		name := v["name"].(string)

		for _, ruleConfig := range v["rewrite_rule"].([]any) {
			r := ruleConfig.(map[string]any)
			conditions := make([]applicationgateways.ApplicationGatewayRewriteRuleCondition, 0)
			requestConfigurations := make([]applicationgateways.ApplicationGatewayHeaderConfiguration, 0)
			responseConfigurations := make([]applicationgateways.ApplicationGatewayHeaderConfiguration, 0)
			urlConfiguration := applicationgateways.ApplicationGatewayURLConfiguration{}

			rule := applicationgateways.ApplicationGatewayRewriteRule{
				Name:         pointer.To(r["name"].(string)),
				RuleSequence: pointer.To(int64(r["rule_sequence"].(int))),
			}

			for _, rawCondition := range r["condition"].([]any) {
				c := rawCondition.(map[string]any)
				condition := applicationgateways.ApplicationGatewayRewriteRuleCondition{
					Variable:   pointer.To(c["variable"].(string)),
					Pattern:    pointer.To(c["pattern"].(string)),
					IgnoreCase: pointer.To(c["ignore_case"].(bool)),
					Negate:     pointer.To(c["negate"].(bool)),
				}
				conditions = append(conditions, condition)
			}
			rule.Conditions = &conditions

			for _, rawConfig := range r["request_header_configuration"].([]any) {
				c := rawConfig.(map[string]any)
				config := applicationgateways.ApplicationGatewayHeaderConfiguration{
					HeaderName:  pointer.To(c["header_name"].(string)),
					HeaderValue: pointer.To(c["header_value"].(string)),
				}
				requestConfigurations = append(requestConfigurations, config)
			}

			for _, rawConfig := range r["response_header_configuration"].([]any) {
				c := rawConfig.(map[string]any)
				config := applicationgateways.ApplicationGatewayHeaderConfiguration{
					HeaderName:  pointer.To(c["header_name"].(string)),
					HeaderValue: pointer.To(c["header_value"].(string)),
				}
				responseConfigurations = append(responseConfigurations, config)
			}

			for _, rawConfig := range r["url"].([]any) {
				c := rawConfig.(map[string]any)
				if c["path"] == nil && c["query_string"] == nil {
					return nil, errors.New("at least one of `path` or `query_string` must be set")
				}
				components := ""
				if c["components"] != nil {
					components = c["components"].(string)
				}
				if c["path"] != nil && components != "query_string_only" {
					urlConfiguration.ModifiedPath = pointer.To(c["path"].(string))
				}
				if c["query_string"] != nil && components != "path_only" {
					urlConfiguration.ModifiedQueryString = pointer.To(c["query_string"].(string))
				}
				if c["reroute"] != nil {
					urlConfiguration.Reroute = pointer.To(c["reroute"].(bool))
				}
			}

			rule.ActionSet = &applicationgateways.ApplicationGatewayRewriteRuleActionSet{
				RequestHeaderConfigurations:  &requestConfigurations,
				ResponseHeaderConfigurations: &responseConfigurations,
			}

			if len(r["url"].([]any)) > 0 {
				rule.ActionSet.UrlConfiguration = &urlConfiguration
			}

			rules = append(rules, rule)
		}

		ruleSet := applicationgateways.ApplicationGatewayRewriteRuleSet{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRewriteRuleSetPropertiesFormat{
				RewriteRules: &rules,
			},
		}

		ruleSets = append(ruleSets, ruleSet)
	}

	return &ruleSets, nil
}

func expandApplicationGatewayRedirectConfigurations(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayRedirectConfiguration, error) {
	vs := d.Get("redirect_configuration").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayRedirectConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		redirectType := v["redirect_type"].(string)
		targetListenerName := v["target_listener_name"].(string)
		targetUrl := v["target_url"].(string)
		includePath := v["include_path"].(bool)
		includeQueryString := v["include_query_string"].(bool)

		output := applicationgateways.ApplicationGatewayRedirectConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRedirectConfigurationPropertiesFormat{
				RedirectType:       pointer.ToEnum[applicationgateways.ApplicationGatewayRedirectType](redirectType),
				IncludeQueryString: pointer.To(includeQueryString),
				IncludePath:        pointer.To(includePath),
			},
		}

		if targetListenerName == "" && targetUrl == "" {
			return nil, errors.New("set either `target_listener_name` or `target_url`")
		}

		if targetListenerName != "" && targetUrl != "" {
			return nil, errors.New("conflict between `target_listener_name` and `target_url` (redirection is either to URL or target listener)")
		}

		if targetListenerName != "" {
			targetListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, targetListenerName)
			output.Properties.TargetListener = &applicationgateways.SubResource{
				Id: pointer.To(targetListenerID),
			}
		}

		if targetUrl != "" {
			output.Properties.TargetURL = pointer.To(targetUrl)
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandApplicationGatewayAutoscaleConfiguration(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewayAutoscaleConfiguration {
	vs := d.Get("autoscale_configuration").([]any)
	if len(vs) == 0 {
		return nil
	}
	v := vs[0].(map[string]any)

	minCapacity := int64(v["min_capacity"].(int))
	maxCapacity := int64(v["max_capacity"].(int))

	configuration := applicationgateways.ApplicationGatewayAutoscaleConfiguration{
		MinCapacity: minCapacity,
	}

	if maxCapacity != 0 {
		configuration.MaxCapacity = pointer.To(maxCapacity)
	}

	return &configuration
}

func expandApplicationGatewaySku(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewaySku {
	vs := d.Get("sku").([]any)
	v := vs[0].(map[string]any)

	name := v["name"].(string)
	tier := v["tier"].(string)
	capacity := int64(v["capacity"].(int))

	sku := applicationgateways.ApplicationGatewaySku{
		Name: pointer.ToEnum[applicationgateways.ApplicationGatewaySkuName](name),
		Tier: pointer.ToEnum[applicationgateways.ApplicationGatewayTier](tier),
	}

	if capacity != 0 {
		sku.Capacity = pointer.To(capacity)
	}

	return &sku
}

func expandApplicationGatewaySslCertificates(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewaySslCertificate, error) {
	vs := d.Get("ssl_certificate").(*schema.Set).List()
	results := make([]applicationgateways.ApplicationGatewaySslCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		data := v["data"].(string)
		password := v["password"].(string)
		kvsid := v["key_vault_secret_id"].(string)
		cert := v["public_cert_data"].(string)

		output := applicationgateways.ApplicationGatewaySslCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewaySslCertificatePropertiesFormat{},
		}

		//nolint:gocritic
		if data != "" && kvsid != "" {
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		} else if data != "" {
			// data must be base64 encoded
			output.Properties.Data = pointer.To(base64.EncodeIfNot(data))

			output.Properties.Password = pointer.To(password)
		} else if kvsid != "" {
			if password != "" {
				return nil, fmt.Errorf("only one of `key_vault_secret_id` or `password` must be specified for the `ssl_certificate` block %q", name)
			}

			output.Properties.KeyVaultSecretId = pointer.To(kvsid)
		} else if cert != "" {
			output.Properties.PublicCertData = pointer.To(cert)
		} else {
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandApplicationGatewayTrustedClientCertificates(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayTrustedClientCertificate, error) {
	vs := d.Get("trusted_client_certificate").([]any)
	results := make([]applicationgateways.ApplicationGatewayTrustedClientCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)
		data := v["data"].(string)

		output := applicationgateways.ApplicationGatewayTrustedClientCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayTrustedClientCertificatePropertiesFormat{},
		}

		if data != "" {
			// data must be base64 encoded
			output.Properties.Data = pointer.To(base64.EncodeIfNot(data))
		} else {
			return nil, fmt.Errorf("`data` must be specified for the `trusted_client_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandApplicationGatewaySslProfiles(d *pluginsdk.ResourceData, gatewayID string) *[]applicationgateways.ApplicationGatewaySslProfile {
	vs := d.Get("ssl_profile").([]any)
	results := make([]applicationgateways.ApplicationGatewaySslProfile, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)

		verifyClientCertificateRevocation := applicationgateways.ApplicationGatewayClientRevocationOptionsNone
		if v["verify_client_certificate_revocation"].(string) != "" {
			verifyClientCertificateRevocation = applicationgateways.ApplicationGatewayClientRevocationOptions(v["verify_client_certificate_revocation"].(string))
		}

		output := applicationgateways.ApplicationGatewaySslProfile{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewaySslProfilePropertiesFormat{
				ClientAuthConfiguration: &applicationgateways.ApplicationGatewayClientAuthConfiguration{
					VerifyClientCertIssuerDN: pointer.To(v["verify_client_certificate_issuer_dn"].(bool)),
					VerifyClientAuthMode:     pointer.ToEnum[applicationgateways.ApplicationGatewayClientAuthVerificationModes](v["client_authentication_mode"].(string)),
					VerifyClientRevocation:   pointer.To(verifyClientCertificateRevocation),
				},
			},
		}

		if v["trusted_client_certificate_names"] != nil {
			clientCerts := v["trusted_client_certificate_names"].([]any)
			clientCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawClientCert := range clientCerts {
				clientCertName := rawClientCert
				clientCertID := fmt.Sprintf("%s/trustedClientCertificates/%s", gatewayID, clientCertName)
				clientCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(clientCertID),
				}
				clientCertSubResources = append(clientCertSubResources, clientCertSubResource)
			}
			output.Properties.TrustedClientCertificates = &clientCertSubResources
		}

		sslPolicy := v["ssl_policy"].([]any)
		if len(sslPolicy) > 0 {
			output.Properties.SslPolicy = expandApplicationGatewaySslPolicy(sslPolicy)
		} else {
			output.Properties.SslPolicy = nil
		}
		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayURLPathMaps(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayURLPathMap, error) {
	vs := d.Get("url_path_map").([]any)
	results := make([]applicationgateways.ApplicationGatewayURLPathMap, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)

		name := v["name"].(string)

		pathRules := make([]applicationgateways.ApplicationGatewayPathRule, 0)
		for _, ruleConfig := range v["path_rule"].([]any) {
			ruleConfigMap := ruleConfig.(map[string]any)

			ruleName := ruleConfigMap["name"].(string)
			backendAddressPoolName := ruleConfigMap["backend_address_pool_name"].(string)
			backendHTTPSettingsName := ruleConfigMap["backend_http_settings_name"].(string)
			redirectConfigurationName := ruleConfigMap["redirect_configuration_name"].(string)
			firewallPolicyID := ruleConfigMap["firewall_policy_id"].(string)

			rulePaths := make([]string, 0)
			for _, rulePath := range ruleConfigMap["paths"].([]any) {
				p, ok := rulePath.(string)
				if ok {
					rulePaths = append(rulePaths, p)
				}
			}

			rule := applicationgateways.ApplicationGatewayPathRule{
				Name: pointer.To(ruleName),
				Properties: &applicationgateways.ApplicationGatewayPathRulePropertiesFormat{
					Paths: &rulePaths,
				},
			}

			if backendAddressPoolName != "" && redirectConfigurationName != "" {
				return nil, errors.New("conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
			}

			if backendHTTPSettingsName != "" && redirectConfigurationName != "" {
				return nil, errors.New("conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
			}

			if backendAddressPoolName != "" {
				backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
				rule.Properties.BackendAddressPool = &applicationgateways.SubResource{
					Id: pointer.To(backendAddressPoolID),
				}
			}

			if backendHTTPSettingsName != "" {
				backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
				rule.Properties.BackendHTTPSettings = &applicationgateways.SubResource{
					Id: pointer.To(backendHTTPSettingsID),
				}
			}

			if redirectConfigurationName != "" {
				redirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigurationName)
				rule.Properties.RedirectConfiguration = &applicationgateways.SubResource{
					Id: pointer.To(redirectConfigurationID),
				}
			}

			if rewriteRuleSetName := ruleConfigMap["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
				rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
				rule.Properties.RewriteRuleSet = &applicationgateways.SubResource{
					Id: pointer.To(rewriteRuleSetID),
				}
			}

			if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
				rule.Properties.FirewallPolicy = &applicationgateways.SubResource{
					Id: pointer.To(firewallPolicyID),
				}
			}

			pathRules = append(pathRules, rule)
		}

		output := applicationgateways.ApplicationGatewayURLPathMap{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayURLPathMapPropertiesFormat{
				PathRules: &pathRules,
			},
		}

		defaultBackendAddressPoolName := v["default_backend_address_pool_name"].(string)
		defaultBackendHTTPSettingsName := v["default_backend_http_settings_name"].(string)
		defaultRedirectConfigurationName := v["default_redirect_configuration_name"].(string)

		if defaultBackendAddressPoolName == "" && defaultBackendHTTPSettingsName == "" && defaultRedirectConfigurationName == "" {
			return nil, fmt.Errorf("both the `default_backend_address_pool_name` and `default_backend_http_settings_name` or `default_redirect_configuration_name` must be specified")
		}

		if defaultBackendAddressPoolName != "" && defaultRedirectConfigurationName != "" {
			return nil, errors.New("conflict between `default_backend_address_pool_name` and `default_redirect_configuration_name` (back-end pool not applicable when redirection specified)")
		}

		if defaultBackendHTTPSettingsName != "" && defaultRedirectConfigurationName != "" {
			return nil, errors.New("conflict between `default_backend_http_settings_name` and `default_redirect_configuration_name` (back-end settings not applicable when redirection specified)")
		}

		if defaultBackendAddressPoolName != "" {
			defaultBackendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, defaultBackendAddressPoolName)
			output.Properties.DefaultBackendAddressPool = &applicationgateways.SubResource{
				Id: pointer.To(defaultBackendAddressPoolID),
			}
		}

		if defaultBackendHTTPSettingsName != "" {
			defaultBackendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, defaultBackendHTTPSettingsName)
			output.Properties.DefaultBackendHTTPSettings = &applicationgateways.SubResource{
				Id: pointer.To(defaultBackendHTTPSettingsID),
			}
		}

		if defaultRedirectConfigurationName != "" {
			defaultRedirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, defaultRedirectConfigurationName)
			output.Properties.DefaultRedirectConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(defaultRedirectConfigurationID),
			}
		}

		if defaultRewriteRuleSetName := v["default_rewrite_rule_set_name"].(string); defaultRewriteRuleSetName != "" {
			defaultRewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, defaultRewriteRuleSetName)
			output.Properties.DefaultRewriteRuleSet = &applicationgateways.SubResource{
				Id: pointer.To(defaultRewriteRuleSetID),
			}
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandApplicationGatewayWafConfig(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration {
	vs := d.Get("waf_configuration").([]any)
	if len(vs) == 0 || vs[0] == nil {
		return nil
	}
	v := vs[0].(map[string]any)

	enabled := v["enabled"].(bool)
	mode := v["firewall_mode"].(string)
	ruleSetType := v["rule_set_type"].(string)
	ruleSetVersion := v["rule_set_version"].(string)
	fileUploadLimitInMb := v["file_upload_limit_mb"].(int)
	requestBodyCheck := v["request_body_check"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_kb"].(int)

	return &applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:                enabled,
		FirewallMode:           applicationgateways.ApplicationGatewayFirewallMode(mode),
		RuleSetType:            ruleSetType,
		RuleSetVersion:         ruleSetVersion,
		FileUploadLimitInMb:    pointer.To(int64(fileUploadLimitInMb)),
		RequestBodyCheck:       pointer.To(requestBodyCheck),
		MaxRequestBodySizeInKb: pointer.To(int64(maxRequestBodySizeInKb)),
		DisabledRuleGroups:     expandApplicationGatewayFirewallDisabledRuleGroup(v["disabled_rule_group"].([]any)),
		Exclusions:             expandApplicationGatewayFirewallExclusion(v["exclusion"].([]any)),
	}
}

func expandApplicationGatewayFirewallDisabledRuleGroup(d []any) *[]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup {
	if len(d) == 0 {
		return nil
	}

	disabledRuleGroups := make([]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup, 0)
	for _, disabledRuleGroup := range d {
		disabledRuleGroupMap := disabledRuleGroup.(map[string]any)

		ruleGroupName := disabledRuleGroupMap["rule_group_name"].(string)

		ruleGroup := applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup{
			RuleGroupName: ruleGroupName,
		}

		rules := make([]int64, 0)
		for _, rule := range disabledRuleGroupMap["rules"].([]any) {
			rules = append(rules, int64(rule.(int)))
		}

		if len(rules) > 0 {
			ruleGroup.Rules = &rules
		}

		disabledRuleGroups = append(disabledRuleGroups, ruleGroup)
	}
	return &disabledRuleGroups
}

func expandApplicationGatewayFirewallExclusion(d []any) *[]applicationgateways.ApplicationGatewayFirewallExclusion {
	if len(d) == 0 {
		return nil
	}

	exclusions := make([]applicationgateways.ApplicationGatewayFirewallExclusion, 0)
	for _, exclusion := range d {
		exclusionMap := exclusion.(map[string]any)

		matchVariable := exclusionMap["match_variable"].(string)
		selectorMatchOperator := exclusionMap["selector_match_operator"].(string)
		selector := exclusionMap["selector"].(string)

		exclusionList := applicationgateways.ApplicationGatewayFirewallExclusion{
			MatchVariable:         matchVariable,
			SelectorMatchOperator: selectorMatchOperator,
			Selector:              selector,
		}

		exclusions = append(exclusions, exclusionList)
	}

	return &exclusions
}

func expandApplicationGatewayCustomErrorConfigurations(vs []any) *[]applicationgateways.ApplicationGatewayCustomError {
	results := make([]applicationgateways.ApplicationGatewayCustomError, 0)

	for _, raw := range vs {
		v := raw.(map[string]any)
		statusCode := v["status_code"].(string)
		customErrorPageUrl := v["custom_error_page_url"].(string)

		output := applicationgateways.ApplicationGatewayCustomError{
			StatusCode:         pointer.ToEnum[applicationgateways.ApplicationGatewayCustomErrorStatusCode](statusCode),
			CustomErrorPageURL: pointer.To(customErrorPageUrl),
		}
		results = append(results, output)
	}

	return &results
}
