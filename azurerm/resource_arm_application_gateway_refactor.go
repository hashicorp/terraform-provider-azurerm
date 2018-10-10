package azurerm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func expandApplicationGatewaySku(d *schema.ResourceData) *network.ApplicationGatewaySku {
	skuSet := d.Get("sku").([]interface{})
	sku := skuSet[0].(map[string]interface{})

	name := sku["name"].(string)
	tier := sku["tier"].(string)
	capacity := int32(sku["capacity"].(int))

	return &network.ApplicationGatewaySku{
		Name:     network.ApplicationGatewaySkuName(name),
		Tier:     network.ApplicationGatewayTier(tier),
		Capacity: utils.Int32(capacity),
	}
}

func expandApplicationGatewayWafConfig(d *schema.ResourceData) *network.ApplicationGatewayWebApplicationFirewallConfiguration {
	wafSet := d.Get("waf_configuration").(*schema.Set).List()
	waf := wafSet[0].(map[string]interface{})

	enabled := waf["enabled"].(bool)
	mode := waf["firewall_mode"].(string)
	rulesettype := waf["rule_set_type"].(string)
	rulesetversion := waf["rule_set_version"].(string)

	return &network.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:        &enabled,
		FirewallMode:   network.ApplicationGatewayFirewallMode(mode),
		RuleSetType:    &rulesettype,
		RuleSetVersion: &rulesetversion,
	}
}

func expandApplicationGatewayFrontendPorts(d *schema.ResourceData) *[]network.ApplicationGatewayFrontendPort {
	configs := d.Get("frontend_port").([]interface{})
	frontendPorts := make([]network.ApplicationGatewayFrontendPort, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		port := int32(data["port"].(int))

		portConfig := network.ApplicationGatewayFrontendPort{
			Name: &name,
			ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
				Port: &port,
			},
		}
		frontendPorts = append(frontendPorts, portConfig)
	}

	return &frontendPorts
}

func expandApplicationGatewayFrontendIPConfigurations(d *schema.ResourceData) *[]network.ApplicationGatewayFrontendIPConfiguration {
	configs := d.Get("frontend_ip_configuration").([]interface{})
	frontEndConfigs := make([]network.ApplicationGatewayFrontendIPConfiguration, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		properties := network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{}

		if v := data["subnet_id"].(string); v != "" {
			properties.Subnet = &network.SubResource{
				ID: &v,
			}
		}

		if v := data["private_ip_address_allocation"].(string); v != "" {
			properties.PrivateIPAllocationMethod = network.IPAllocationMethod(v)
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &network.SubResource{
				ID: &v,
			}
		}

		name := data["name"].(string)
		frontEndConfig := network.ApplicationGatewayFrontendIPConfiguration{
			Name: &name,
			ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &properties,
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs
}

func expandApplicationGatewayProbes(d *schema.ResourceData) *[]network.ApplicationGatewayProbe {
	configs := d.Get("probe").([]interface{})
	backendSettings := make([]network.ApplicationGatewayProbe, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		protocol := data["protocol"].(string)
		probePath := data["path"].(string)
		host := data["host"].(string)
		interval := int32(data["interval"].(int))
		timeout := int32(data["timeout"].(int))
		unhealthyThreshold := int32(data["unhealthy_threshold"].(int))
		minServers := int32(data["minimum_servers"].(int))

		setting := network.ApplicationGatewayProbe{
			Name: &name,
			ApplicationGatewayProbePropertiesFormat: &network.ApplicationGatewayProbePropertiesFormat{
				Protocol:           network.ApplicationGatewayProtocol(protocol),
				Path:               &probePath,
				Host:               &host,
				Interval:           &interval,
				Timeout:            &timeout,
				UnhealthyThreshold: &unhealthyThreshold,
				MinServers:         &minServers,
			},
		}

		matchConfigs := data["match"].([]interface{})
		if len(matchConfigs) > 0 {
			match := matchConfigs[0].(map[string]interface{})
			matchBody := match["body"].(string)

			statusCodes := make([]string, 0)
			for _, statusCode := range match["status_code"].([]interface{}) {
				statusCodes = append(statusCodes, statusCode.(string))
			}

			setting.ApplicationGatewayProbePropertiesFormat.Match = &network.ApplicationGatewayProbeHealthResponseMatch{
				Body:        &matchBody,
				StatusCodes: &statusCodes,
			}
		}

		backendSettings = append(backendSettings, setting)
	}

	return &backendSettings
}

func expandApplicationGatewayRequestRoutingRules(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayRequestRoutingRule {
	configs := d.Get("request_routing_rule").([]interface{})
	rules := make([]network.ApplicationGatewayRequestRoutingRule, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		ruleType := data["rule_type"].(string)
		httpListenerName := data["http_listener_name"].(string)
		httpListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, httpListenerName)

		rule := network.ApplicationGatewayRequestRoutingRule{
			Name: &name,
			ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
				RuleType: network.ApplicationGatewayRequestRoutingRuleType(ruleType),
				HTTPListener: &network.SubResource{
					ID: &httpListenerID,
				},
			},
		}

		if backendAddressPoolName := data["backend_address_pool_name"].(string); backendAddressPoolName != "" {
			backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendAddressPool = &network.SubResource{
				ID: &backendAddressPoolID,
			}
		}

		if backendHTTPSettingsName := data["backend_http_settings_name"].(string); backendHTTPSettingsName != "" {
			backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
				ID: &backendHTTPSettingsID,
			}
		}

		if urlPathMapName := data["url_path_map_name"].(string); urlPathMapName != "" {
			urlPathMapID := fmt.Sprintf("%s/urlPathMaps/%s", gatewayID, urlPathMapName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.URLPathMap = &network.SubResource{
				ID: &urlPathMapID,
			}
		}

		rules = append(rules, rule)
	}

	return &rules
}

func expandApplicationGatewayURLPathMaps(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayURLPathMap {
	configs := d.Get("url_path_map").([]interface{})
	pathMaps := make([]network.ApplicationGatewayURLPathMap, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		defaultBackendAddressPoolName := data["default_backend_address_pool_name"].(string)
		defaultBackendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, defaultBackendAddressPoolName)
		defaultBackendHTTPSettingsName := data["default_backend_http_settings_name"].(string)
		defaultBackendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, defaultBackendHTTPSettingsName)

		pathRules := []network.ApplicationGatewayPathRule{}
		for _, ruleConfig := range data["path_rule"].([]interface{}) {
			ruleConfigMap := ruleConfig.(map[string]interface{})

			ruleName := ruleConfigMap["name"].(string)

			rulePaths := []string{}
			for _, rulePath := range ruleConfigMap["paths"].([]interface{}) {
				rulePaths = append(rulePaths, rulePath.(string))
			}

			rule := network.ApplicationGatewayPathRule{
				Name: &ruleName,
				ApplicationGatewayPathRulePropertiesFormat: &network.ApplicationGatewayPathRulePropertiesFormat{
					Paths: &rulePaths,
				},
			}

			if backendAddressPoolName := ruleConfigMap["backend_address_pool_name"].(string); backendAddressPoolName != "" {
				backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendAddressPool = &network.SubResource{
					ID: &backendAddressPoolID,
				}
			}

			if backendHTTPSettingsName := ruleConfigMap["backend_http_settings_name"].(string); backendHTTPSettingsName != "" {
				backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
					ID: &backendHTTPSettingsID,
				}
			}

			pathRules = append(pathRules, rule)
		}

		pathMap := network.ApplicationGatewayURLPathMap{
			Name: &name,
			ApplicationGatewayURLPathMapPropertiesFormat: &network.ApplicationGatewayURLPathMapPropertiesFormat{
				DefaultBackendAddressPool: &network.SubResource{
					ID: &defaultBackendAddressPoolID,
				},
				DefaultBackendHTTPSettings: &network.SubResource{
					ID: &defaultBackendHTTPSettingsID,
				},
				PathRules: &pathRules,
			},
		}

		pathMaps = append(pathMaps, pathMap)
	}

	return &pathMaps
}

func expandApplicationGatewaySslCertificates(d *schema.ResourceData) *[]network.ApplicationGatewaySslCertificate {
	configs := d.Get("ssl_certificate").([]interface{})
	sslCerts := make([]network.ApplicationGatewaySslCertificate, 0)

	for _, configRaw := range configs {
		raw := configRaw.(map[string]interface{})

		name := raw["name"].(string)
		data := raw["data"].(string)
		password := raw["password"].(string)

		// data must be base64 encoded
		data = base64Encode(data)

		cert := network.ApplicationGatewaySslCertificate{
			Name: &name,
			ApplicationGatewaySslCertificatePropertiesFormat: &network.ApplicationGatewaySslCertificatePropertiesFormat{
				Data:     &data,
				Password: &password,
			},
		}

		sslCerts = append(sslCerts, cert)
	}

	return &sslCerts
}

func flattenApplicationGatewaySku(sku *network.ApplicationGatewaySku) []interface{} {
	result := make(map[string]interface{})

	result["name"] = string(sku.Name)
	result["tier"] = string(sku.Tier)
	result["capacity"] = int(*sku.Capacity)

	return []interface{}{result}
}

func flattenApplicationGatewayWafConfig(waf *network.ApplicationGatewayWebApplicationFirewallConfiguration) []interface{} {
	result := make(map[string]interface{})

	result["enabled"] = *waf.Enabled
	result["firewall_mode"] = string(waf.FirewallMode)
	result["rule_set_type"] = waf.RuleSetType
	result["rule_set_version"] = waf.RuleSetVersion

	return []interface{}{result}
}

func flattenApplicationGatewayFrontendPorts(portConfigs *[]network.ApplicationGatewayFrontendPort) []interface{} {
	result := make([]interface{}, 0)

	if configs := portConfigs; configs != nil {
		for _, config := range *configs {
			port := map[string]interface{}{
				"id":   *config.ID,
				"name": *config.Name,
				"port": int(*config.ApplicationGatewayFrontendPortPropertiesFormat.Port),
			}

			result = append(result, port)
		}
	}

	return result
}

func flattenApplicationGatewayFrontendIPConfigurations(ipConfigs *[]network.ApplicationGatewayFrontendIPConfiguration) []interface{} {
	result := make([]interface{}, 0)

	if configs := ipConfigs; configs != nil {
		for _, config := range *ipConfigs {
			ipConfig := make(map[string]interface{})
			ipConfig["id"] = *config.ID
			ipConfig["name"] = *config.Name

			if props := config.ApplicationGatewayFrontendIPConfigurationPropertiesFormat; props != nil {
				if props.PrivateIPAllocationMethod != "" {
					ipConfig["private_ip_address_allocation"] = props.PrivateIPAllocationMethod
				}

				if props.Subnet != nil {
					ipConfig["subnet_id"] = *props.Subnet.ID
				}

				if props.PrivateIPAddress != nil {
					ipConfig["private_ip_address"] = *props.PrivateIPAddress
				}

				if props.PublicIPAddress != nil {
					ipConfig["public_ip_address_id"] = *props.PublicIPAddress.ID
				}
			}

			result = append(result, ipConfig)
		}
	}

	return result
}

func flattenApplicationGatewayProbes(input *[]network.ApplicationGatewayProbe) []interface{} {
	result := make([]interface{}, 0)

	if probes := input; probes != nil {
		for _, config := range *probes {
			settings := map[string]interface{}{
				"id":   *config.ID,
				"name": *config.Name,
			}

			if props := config.ApplicationGatewayProbePropertiesFormat; props != nil {
				settings["protocol"] = string(props.Protocol)

				if host := props.Host; host != nil {
					settings["host"] = *host
				}

				if path := props.Path; path != nil {
					settings["path"] = *path
				}

				if interval := props.Interval; interval != nil {
					settings["interval"] = int(*interval)
				}

				if timeout := props.Timeout; timeout != nil {
					settings["timeout"] = int(*timeout)
				}

				if threshold := props.UnhealthyThreshold; threshold != nil {
					settings["unhealthy_threshold"] = int(*threshold)
				}

				if minServers := props.MinServers; minServers != nil {
					settings["minimum_servers"] = int(*minServers)
				}

				if match := props.Match; match != nil {
					matchConfig := map[string]interface{}{}
					if body := match.Body; body != nil {
						matchConfig["body"] = *body
					}

					statusCodes := make([]interface{}, 0)
					if match.StatusCodes != nil {
						for _, status := range *match.StatusCodes {
							statusCodes = append(statusCodes, status)
						}
						matchConfig["status_code"] = statusCodes
					}
					settings["match"] = matchConfig
				}
			}

			result = append(result, settings)
		}
	}

	return result
}

func flattenApplicationGatewayRequestRoutingRules(input *[]network.ApplicationGatewayRequestRoutingRule) ([]interface{}, error) {
	result := make([]interface{}, 0)

	if rules := input; rules != nil {
		for _, config := range *rules {

			if props := config.ApplicationGatewayRequestRoutingRulePropertiesFormat; props != nil {
				httpListenerName := strings.Split(*props.HTTPListener.ID, "/")[len(strings.Split(*props.HTTPListener.ID, "/"))-1]
				listener := map[string]interface{}{
					"id":                 *config.ID,
					"name":               *config.Name,
					"rule_type":          string(props.RuleType),
					"http_listener_id":   *props.HTTPListener.ID,
					"http_listener_name": httpListenerName,
				}

				if pool := props.BackendAddressPool; pool != nil {
					backendAddressPoolName := strings.Split(*pool.ID, "/")[len(strings.Split(*pool.ID, "/"))-1]
					listener["backend_address_pool_name"] = backendAddressPoolName
					listener["backend_address_pool_id"] = *pool.ID
				}

				if settings := props.BackendHTTPSettings; settings != nil {
					backendHTTPSettingsName := strings.Split(*settings.ID, "/")[len(strings.Split(*settings.ID, "/"))-1]
					listener["backend_http_settings_name"] = backendHTTPSettingsName
					listener["backend_http_settings_id"] = *settings.ID
				}

				if pathMap := props.URLPathMap; pathMap != nil {
					urlPathMapName := strings.Split(*pathMap.ID, "/")[len(strings.Split(*pathMap.ID, "/"))-1]
					listener["url_path_map_name"] = urlPathMapName
					listener["url_path_map_id"] = *pathMap.ID
				}

				result = append(result, listener)
			}
		}
	}

	return result, nil
}

func flattenApplicationGatewayURLPathMaps(input *[]network.ApplicationGatewayURLPathMap) ([]interface{}, error) {
	result := make([]interface{}, 0)

	if pathMaps := input; pathMaps != nil {
		for _, config := range *pathMaps {
			pathMap := map[string]interface{}{
				"id":   *config.ID,
				"name": *config.Name,
			}

			if props := config.ApplicationGatewayURLPathMapPropertiesFormat; props != nil {
				if backendPool := props.DefaultBackendAddressPool; backendPool != nil {
					backendAddressPoolName := strings.Split(*backendPool.ID, "/")[len(strings.Split(*backendPool.ID, "/"))-1]
					pathMap["default_backend_address_pool_name"] = backendAddressPoolName
					pathMap["default_backend_address_pool_id"] = *backendPool.ID
				}

				if settings := props.DefaultBackendHTTPSettings; settings != nil {
					backendHTTPSettingsName := strings.Split(*settings.ID, "/")[len(strings.Split(*settings.ID, "/"))-1]
					pathMap["default_backend_http_settings_name"] = backendHTTPSettingsName
					pathMap["default_backend_http_settings_id"] = *settings.ID
				}

				pathRules := make([]interface{}, 0)
				if rules := props.PathRules; rules != nil {
					for _, pathRuleConfig := range *rules {
						rule := map[string]interface{}{
							"id":   *pathRuleConfig.ID,
							"name": *pathRuleConfig.Name,
						}

						if ruleProps := pathRuleConfig.ApplicationGatewayPathRulePropertiesFormat; props != nil {
							if pool := ruleProps.BackendAddressPool; pool != nil {
								backendAddressPoolName2 := strings.Split(*pool.ID, "/")[len(strings.Split(*pool.ID, "/"))-1]
								rule["backend_address_pool_name"] = backendAddressPoolName2
								rule["backend_address_pool_id"] = *pool.ID
							}

							if backend := ruleProps.BackendHTTPSettings; backend != nil {
								backendHTTPSettingsName2 := strings.Split(*backend.ID, "/")[len(strings.Split(*backend.ID, "/"))-1]
								rule["backend_http_settings_name"] = backendHTTPSettingsName2
								rule["backend_http_settings_id"] = *backend.ID
							}

							pathOutputs := make([]interface{}, 0)
							if paths := ruleProps.Paths; paths != nil {
								for _, rulePath := range *paths {
									pathOutputs = append(pathOutputs, rulePath)
								}
							}
							rule["paths"] = pathOutputs
						}

						pathRules = append(pathRules, rule)
					}
					pathMap["path_rule"] = pathRules
				}
			}

			result = append(result, pathMap)
		}
	}

	return result, nil
}

func flattenApplicationGatewaySslCertificates(input *[]network.ApplicationGatewaySslCertificate) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewaySslCertificatePropertiesFormat; props != nil {
			if data := props.PublicCertData; data != nil {
				output["public_cert_data"] = *data
			}
		}

		results = append(results, output)
	}

	return results
}

// TODO: can this be removed?
func hashApplicationGatewayWafConfig(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%t-", m["enabled"].(bool)))
	buf.WriteString(fmt.Sprintf("%s-", m["firewall_mode"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", *m["rule_set_type"].(*string)))
	buf.WriteString(fmt.Sprintf("%s-", *m["rule_set_version"].(*string)))

	return hashcode.String(buf.String())
}

func hashApplicationGatewaySslCertificates(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["public_cert_data"].(string)))

	return hashcode.String(buf.String())
}
