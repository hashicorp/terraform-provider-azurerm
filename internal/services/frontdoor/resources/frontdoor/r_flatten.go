// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
)

func flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings []any, frontDoorId frontdoors.FrontDoorId) *[]any {
	output := make([]any, 0)
	var backendPoolOrder []string
	var frontedEndpointOrder []string
	var routingRulesOrder []string
	var backendPoolLoadBalancingOrder []string
	var backendPoolHealthProbeOrder []string
	if len(backendPools) > 0 {
		flattenendBackendPools, err := flattenFrontDoorBackendPools(expandFrontDoorBackendPools(backendPools, frontDoorId), frontDoorId, make([]any, 0))
		if err == nil {
			for _, ids := range *flattenendBackendPools {
				backendPool := ids.(map[string]any)
				backendPoolOrder = append(backendPoolOrder, backendPool["id"].(string))
			}
		}
	}
	if len(frontendEndpoints) > 0 {
		flattenendfrontendEndpoints, err := flattenFrontEndEndpoints(expandFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorId), frontDoorId, make([]any, 0))
		if err == nil {
			for _, ids := range *flattenendfrontendEndpoints {
				frontendEndPoint := ids.(map[string]any)
				frontedEndpointOrder = append(frontedEndpointOrder, frontendEndPoint["id"].(string))
			}
		}
	}
	if len(routingRules) > 0 {
		var oldBlocks any
		flattenendRoutingRules, err := flattenFrontDoorRoutingRule(expandFrontDoorRoutingRule(routingRules, frontDoorId, nil), oldBlocks, frontDoorId, make([]any, 0))
		if err == nil {
			for _, ids := range *flattenendRoutingRules {
				routingRule := ids.(map[string]any)
				routingRulesOrder = append(routingRulesOrder, routingRule["id"].(string))
			}
		}
	}
	if len(loadBalancingSettings) > 0 {
		flattenendLoadBalancingSettings := flattenFrontDoorLoadBalancingSettingsModel(expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorId), frontDoorId, make([]any, 0))

		if len(flattenendLoadBalancingSettings) > 0 {
			for _, ids := range flattenendLoadBalancingSettings {
				loadBalancingSetting := ids.(map[string]any)
				backendPoolLoadBalancingOrder = append(backendPoolLoadBalancingOrder, loadBalancingSetting["id"].(string))
			}
		}
	}
	if len(healthProbeSettings) > 0 {
		flattenendHealthProbeSettings := flattenFrontDoorHealthProbeSettingsModel(expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorId), frontDoorId, make([]any, 0))

		if len(flattenendHealthProbeSettings) > 0 {
			for _, ids := range flattenendHealthProbeSettings {
				healthProbeSetting := ids.(map[string]any)
				backendPoolHealthProbeOrder = append(backendPoolHealthProbeOrder, healthProbeSetting["id"].(string))
			}
		}
	}

	output = append(output, map[string]any{
		"backend_pool_ids":                backendPoolOrder,
		"frontend_endpoint_ids":           frontedEndpointOrder,
		"routing_rule_ids":                routingRulesOrder,
		"backend_pool_load_balancing_ids": backendPoolLoadBalancingOrder,
		"backend_pool_health_probe_ids":   backendPoolHealthProbeOrder,
	})

	return &output
}

func flattenFrontDoorBackendPools(input *[]frontdoors.BackendPool, frontDoorId frontdoors.FrontDoorId, explicitOrder []any) (*[]any, error) {
	if input == nil {
		return &[]any{}, nil
	}

	output := make([]any, 0)

	if len(explicitOrder) > 0 {
		orderedBackendPools := explicitOrder[0].(map[string]any)
		orderedBackendPoolsIds := orderedBackendPools["backend_pool_ids"].([]any)
		combinedBackendPools, err := combineBackendPools(*input, orderedBackendPoolsIds, frontDoorId)
		if err == nil {
			output = combinedBackendPools
		} else {
			return nil, err
		}
	} else {
		for _, backend := range *input {
			backendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, backendPool)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorBackendPools(input *frontdoors.BackendPool, frontDoorId frontdoors.FrontDoorId) (map[string]any, error) {
	if input == nil {
		return make(map[string]any), nil
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	backend := make([]any, 0)
	healthProbeName := ""
	loadBalancingName := ""
	if props := input.Properties; props != nil {
		backend = flattenFrontDoorBackend(props.Backends)
		if props.HealthProbeSettings != nil && props.HealthProbeSettings.Id != nil {
			name, err := parse.HealthProbeIDInsensitively(*props.HealthProbeSettings.Id)
			if err != nil {
				return nil, err
			}
			healthProbeName = name.HealthProbeSettingName
		}

		if props.LoadBalancingSettings != nil && props.LoadBalancingSettings.Id != nil {
			name, err := parse.LoadBalancingIDInsensitively(*props.LoadBalancingSettings.Id)
			if err != nil {
				return nil, err
			}
			loadBalancingName = name.LoadBalancingSettingName
		}
	}

	output := map[string]any{
		"backend":             backend,
		"health_probe_name":   healthProbeName,
		"id":                  id,
		"load_balancing_name": loadBalancingName,
		"name":                name,
	}

	return output, nil
}

func flattenFrontDoorBackendPoolsSettings(input *frontdoors.BackendPoolsSettings) flattenedBackendPoolSettings {
	if input == nil {
		return flattenedBackendPoolSettings{
			enforceBackendPoolsCertificateNameCheck: true,
			backendPoolsSendReceiveTimeoutSeconds:   60,
		}
	}

	enforceCertificateNameCheck := false
	sendReceiveTimeoutSeconds := 0
	if input.EnforceCertificateNameCheck != nil && *input.EnforceCertificateNameCheck != "" && *input.EnforceCertificateNameCheck == frontdoors.EnforceCertificateNameCheckEnabledStateEnabled {
		enforceCertificateNameCheck = true
	}
	if input.SendRecvTimeoutSeconds != nil {
		sendReceiveTimeoutSeconds = int(*input.SendRecvTimeoutSeconds)
	}

	return flattenedBackendPoolSettings{
		enforceBackendPoolsCertificateNameCheck: enforceCertificateNameCheck,
		backendPoolsSendReceiveTimeoutSeconds:   sendReceiveTimeoutSeconds,
	}
}

func flattenFrontDoorBackend(input *[]frontdoors.Backend) []any {
	if input == nil {
		return make([]any, 0)
	}
	output := make([]any, 0)
	for _, v := range *input {
		result := make(map[string]any)
		if address := v.Address; address != nil {
			result["address"] = *address
		}
		if backendHostHeader := v.BackendHostHeader; backendHostHeader != nil {
			result["host_header"] = *backendHostHeader
		}
		if v.EnabledState != nil {
			result["enabled"] = *v.EnabledState == frontdoors.BackendEnabledStateEnabled
		}
		if httpPort := v.HTTPPort; httpPort != nil {
			result["http_port"] = int(*httpPort)
		}
		if httpsPort := v.HTTPSPort; httpsPort != nil {
			result["https_port"] = int(*httpsPort)
		}
		if priority := v.Priority; priority != nil {
			result["priority"] = int(*priority)
		}
		if weight := v.Weight; weight != nil {
			result["weight"] = int(*weight)
		}
		output = append(output, result)
	}

	return output
}

func flattenFrontEndEndpoints(input *[]frontdoors.FrontendEndpoint, frontDoorId frontdoors.FrontDoorId, explicitOrder []any) (*[]any, error) {
	output := make([]any, 0)
	if input == nil {
		return &output, nil
	}

	if len(explicitOrder) > 0 {
		orderedFrontEnd := explicitOrder[0].(map[string]any)
		orderedFrontEndIds := orderedFrontEnd["frontend_endpoint_ids"].([]any)
		combinedFrontEndEndpoints, err := combineFrontEndEndpoints(*input, orderedFrontEndIds, frontDoorId)
		if err == nil {
			output = combinedFrontEndEndpoints
		} else {
			return nil, err
		}
	} else {
		for _, v := range *input {
			frontendEndpoint, err := flattenSingleFrontEndEndpoints(v, frontDoorId)
			if err == nil {
				output = append(output, frontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontEndEndpoints(input frontdoors.FrontendEndpoint, frontDoorId frontdoors.FrontDoorId) (map[string]any, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = frontdoors.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, *input.Name).ID()
		name = *input.Name
	}
	hostName := ""
	sessionAffinityEnabled := false
	sessionAffinityTlsSeconds := 0
	webApplicationFirewallPolicyLinkId := ""
	if props := input.Properties; props != nil {
		if props.HostName != nil {
			hostName = *props.HostName
		}
		if props.SessionAffinityEnabledState != nil && *props.SessionAffinityEnabledState != "" {
			sessionAffinityEnabled = *props.SessionAffinityEnabledState == frontdoors.SessionAffinityEnabledStateEnabled
		}
		if props.SessionAffinityTtlSeconds != nil {
			sessionAffinityTlsSeconds = int(*props.SessionAffinityTtlSeconds)
		}
		if waf := props.WebApplicationFirewallPolicyLink; waf != nil && waf.Id != nil {
			// rewrite the ID to ensure it's consistent
			parsed, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(*waf.Id)
			if err != nil {
				return nil, err
			}
			webApplicationFirewallPolicyLinkId = parsed.ID()
		}
		// flattenedHttpsConfig := flattenCustomHttpsConfiguration(props)
		// customHTTPSConfiguration = flattenedHttpsConfig.CustomHTTPSConfiguration
		// customHttpsProvisioningEnabled = flattenedHttpsConfig.CustomHTTPSProvisioningEnabled
	}

	output := map[string]any{
		// "custom_https_configuration":        customHTTPSConfiguration,
		// "custom_https_provisioning_enabled": customHttpsProvisioningEnabled,
		"host_name":                    hostName,
		"id":                           id,
		"name":                         name,
		"session_affinity_enabled":     sessionAffinityEnabled,
		"session_affinity_ttl_seconds": sessionAffinityTlsSeconds,
		"web_application_firewall_policy_link_id": webApplicationFirewallPolicyLinkId,
	}

	return output, nil
}

func flattenFrontDoorHealthProbeSettingsModel(input *[]frontdoors.HealthProbeSettingsModel, frontDoorId frontdoors.FrontDoorId, explicitOrder []any) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}

	if len(explicitOrder) > 0 {
		orderedHealthProbeSetting := explicitOrder[0].(map[string]any)
		orderedHealthProbeSettingIds := orderedHealthProbeSetting["backend_pool_health_probe_ids"].([]any)
		output = combineHealthProbeSettingsModel(*input, orderedHealthProbeSettingIds, frontDoorId)
	} else {
		for _, v := range *input {
			healthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&v, frontDoorId)
			output = append(output, healthProbeSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorHealthProbeSettingsModel(input *frontdoors.HealthProbeSettingsModel, frontDoorId frontdoors.FrontDoorId) map[string]any {
	if input == nil {
		return make(map[string]any)
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	enabled := false
	intervalInSeconds := 0
	path := ""
	probeMethod := ""
	protocol := ""

	if properties := input.Properties; properties != nil {
		if properties.IntervalInSeconds != nil {
			intervalInSeconds = int(*properties.IntervalInSeconds)
		}
		if properties.Path != nil {
			path = *properties.Path
		}
		if healthProbeMethod := properties.HealthProbeMethod; healthProbeMethod != nil && *healthProbeMethod != "" {
			// I have to upper this as the frontdoor.GET and frontdoor.HEAD types are uppercased
			// but Azure stores them in the resource as sentence cased (e.g. "Get" and "Head")
			probeMethod = strings.ToUpper(string(*healthProbeMethod))
		}
		if properties.EnabledState != nil && *properties.EnabledState != "" {
			enabled = *properties.EnabledState == frontdoors.HealthProbeEnabledEnabled
		}
		if properties.Protocol != nil {
			protocol = string(*properties.Protocol)
		}
	}

	return map[string]any{
		"enabled":             enabled,
		"id":                  id,
		"name":                name,
		"protocol":            protocol,
		"interval_in_seconds": intervalInSeconds,
		"path":                path,
		"probe_method":        probeMethod,
	}
}

func flattenFrontDoorLoadBalancingSettingsModel(input *[]frontdoors.LoadBalancingSettingsModel, frontDoorId frontdoors.FrontDoorId, explicitOrder []any) []any {
	if input == nil {
		return make([]any, 0)
	}

	output := make([]any, 0)

	if len(explicitOrder) > 0 {
		orderedLoadBalancingSettings := explicitOrder[0].(map[string]any)
		orderedLoadBalancingIds := orderedLoadBalancingSettings["backend_pool_load_balancing_ids"].([]any)
		output = combineLoadBalancingSettingsModel(*input, orderedLoadBalancingIds, frontDoorId)
	} else {
		for _, v := range *input {
			loadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&v, frontDoorId)
			output = append(output, loadBalanceSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorLoadBalancingSettingsModel(input *frontdoors.LoadBalancingSettingsModel, frontDoorId frontdoors.FrontDoorId) map[string]any {
	if input == nil {
		return make(map[string]any)
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	additionalLatencyMilliseconds := 0
	sampleSize := 0
	successfulSamplesRequired := 0
	if properties := input.Properties; properties != nil {
		if properties.AdditionalLatencyMilliseconds != nil {
			additionalLatencyMilliseconds = int(*properties.AdditionalLatencyMilliseconds)
		}
		if properties.SampleSize != nil {
			sampleSize = int(*properties.SampleSize)
		}
		if properties.SuccessfulSamplesRequired != nil {
			successfulSamplesRequired = int(*properties.SuccessfulSamplesRequired)
		}
	}

	return map[string]any{
		"additional_latency_milliseconds": additionalLatencyMilliseconds,
		"id":                              id,
		"name":                            name,
		"sample_size":                     sampleSize,
		"successful_samples_required":     successfulSamplesRequired,
	}
}

func flattenFrontDoorRoutingRule(input *[]frontdoors.RoutingRule, oldBlocks any, frontDoorId frontdoors.FrontDoorId, explicitOrder []any) (*[]any, error) {
	if input == nil {
		return &[]any{}, nil
	}

	output := make([]any, 0)

	if len(explicitOrder) > 0 {
		orderedRule := explicitOrder[0].(map[string]any)
		orderedRoutingRuleIds := orderedRule["routing_rule_ids"].([]any)
		combinedRoutingRules, err := combineRoutingRules(*input, oldBlocks, orderedRoutingRuleIds, frontDoorId)
		if err != nil {
			return nil, err
		}
		output = combinedRoutingRules
	} else {
		for _, v := range *input {
			routingRule, err := flattenSingleFrontDoorRoutingRule(v, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, routingRule)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorRoutingRule(input frontdoors.RoutingRule, oldBlocks any, frontDoorId frontdoors.FrontDoorId) (map[string]any, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, *input.Name).ID()
		name = *input.Name
	}

	acceptedProtocols := make([]string, 0)
	enabled := false
	forwardingConfiguration := make([]any, 0)
	frontEndEndpoints := make([]string, 0)
	patternsToMatch := make([]string, 0)
	redirectConfiguration := make([]any, 0)

	if props := input.Properties; props != nil {
		acceptedProtocols = flattenFrontDoorAcceptedProtocol(props.AcceptedProtocols)
		if props.EnabledState != nil {
			enabled = *props.EnabledState == frontdoors.RoutingRuleEnabledStateEnabled
		}
		forwardConfiguration, err := flattenRoutingRuleForwardingConfiguration(props.RouteConfiguration, oldBlocks)
		if err != nil {
			return nil, fmt.Errorf("flattening `forward_configuration`: %+v", err)
		}

		forwardingConfiguration = *forwardConfiguration
		frontendEndpoints, err := flattenFrontDoorFrontendEndpointsSubResources(props.FrontendEndpoints)
		if err != nil {
			return nil, fmt.Errorf("flattening `frontend_endpoints`: %+v", err)
		}

		frontEndEndpoints = *frontendEndpoints
		if props.PatternsToMatch != nil {
			patternsToMatch = *props.PatternsToMatch
		}
		redirectConfiguration = flattenRoutingRuleRedirectConfiguration(props.RouteConfiguration)
	}

	output := map[string]any{
		"accepted_protocols":       acceptedProtocols,
		"enabled":                  enabled,
		"forwarding_configuration": forwardingConfiguration,
		"frontend_endpoints":       frontEndEndpoints,
		"id":                       id,
		"name":                     name,
		"patterns_to_match":        patternsToMatch,
		"redirect_configuration":   redirectConfiguration,
	}

	return output, nil
}

func flattenRoutingRuleForwardingConfiguration(config frontdoors.RouteConfiguration, oldConfig any) (*[]any, error) {
	v, ok := config.(frontdoors.ForwardingConfiguration)
	if !ok {
		return &[]any{}, nil
	}

	name := ""
	if v.BackendPool != nil && v.BackendPool.Id != nil {
		backendPoolId, err := parse.BackendPoolIDInsensitively(*v.BackendPool.Id)
		if err != nil {
			return nil, err
		}
		name = backendPoolId.Name
	}

	cacheEnabled := false
	cacheQueryParameterStripDirective := string(frontdoors.FrontDoorQueryStripAll)
	cacheUseDynamicCompression := false

	var cacheQueryParameters []any
	var cacheQueryParametersArray []string
	var cacheDuration *string

	if cacheConfiguration := v.CacheConfiguration; cacheConfiguration != nil {
		cacheEnabled = true
		if stripDirective := cacheConfiguration.QueryParameterStripDirective; stripDirective != nil && *stripDirective != "" {
			cacheQueryParameterStripDirective = string(*stripDirective)
		}
		if dynamicCompression := cacheConfiguration.DynamicCompression; dynamicCompression != nil && *dynamicCompression != "" {
			cacheUseDynamicCompression = string(*dynamicCompression) == string(frontdoors.DynamicCompressionEnabledEnabled)
		}
		if queryParameters := cacheConfiguration.QueryParameters; queryParameters != nil {
			cacheQueryParametersArray = strings.Split(*queryParameters, ",")
		}
		if duration := cacheConfiguration.CacheDuration; duration != nil {
			cacheDuration = duration
		}
	} else {
		// if the cache is disabled, use the default values or revert to what they were in the previous plan
		old, ok := oldConfig.([]any)
		if ok {
			for _, oldValue := range old {
				oldVal, ok := oldValue.(map[string]any)
				if ok {
					thisName := oldVal["name"].(string)
					if name == thisName {
						oldConfigs := oldVal["forwarding_configuration"].([]any)
						if len(oldConfigs) > 0 {
							ofc := oldConfigs[0].(map[string]any)
							cacheQueryParameterStripDirective = ofc["cache_query_parameter_strip_directive"].(string)
							cacheUseDynamicCompression = ofc["cache_use_dynamic_compression"].(bool)
							cacheDuration = pointer.To(ofc["cache_duration"].(string))

							cacheQueryParameters = ofc["cache_query_parameters"].([]any)
							for _, p := range cacheQueryParameters {
								cacheQueryParametersArray = append(cacheQueryParametersArray, p.(string))
							}
						}
					}
				}
			}
		}
	}

	forwardingProtocol := ""
	if v.ForwardingProtocol != nil {
		forwardingProtocol = string(*v.ForwardingProtocol)
	}

	return &[]any{
		map[string]any{
			"backend_pool_name":                     name,
			"custom_forwarding_path":                pointer.From(v.CustomForwardingPath),
			"forwarding_protocol":                   forwardingProtocol,
			"cache_enabled":                         cacheEnabled,
			"cache_query_parameter_strip_directive": cacheQueryParameterStripDirective,
			"cache_use_dynamic_compression":         cacheUseDynamicCompression,
			"cache_query_parameters":                cacheQueryParametersArray,
			"cache_duration":                        cacheDuration,
		},
	}, nil
}

func flattenRoutingRuleRedirectConfiguration(config frontdoors.RouteConfiguration) []any {
	v, ok := config.(frontdoors.RedirectConfiguration)
	if !ok {
		return []any{}
	}

	redirectProtocol := ""
	if v.RedirectProtocol != nil {
		redirectProtocol = string(*v.RedirectProtocol)
	}

	redirectType := ""
	if v.RedirectType != nil {
		redirectType = string(*v.RedirectType)
	}

	return []any{
		map[string]any{
			"custom_host":         pointer.From(v.CustomHost),
			"custom_fragment":     pointer.From(v.CustomFragment),
			"custom_query_string": pointer.From(v.CustomQueryString),
			"custom_path":         pointer.From(v.CustomPath),
			"redirect_protocol":   redirectProtocol,
			"redirect_type":       redirectType,
		},
	}
}

func flattenFrontDoorAcceptedProtocol(input *[]frontdoors.FrontDoorProtocol) []string {
	if input == nil {
		return make([]string, 0)
	}

	output := make([]string, 0)

	for _, p := range *input {
		output = append(output, string(p))
	}

	return output
}

func flattenFrontDoorFrontendEndpointsSubResources(input *[]frontdoors.SubResource) (*[]string, error) {
	output := make([]string, 0)
	if input == nil {
		return &output, nil
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		id, err := frontdoors.ParseFrontendEndpointIDInsensitively(*v.Id)
		if err != nil {
			return nil, err
		}
		output = append(output, id.FrontendEndpointName)
	}

	return &output, nil
}

func combineBackendPools(allPools []frontdoors.BackendPool, orderedIds []any, frontDoorId frontdoors.FrontDoorId) ([]any, error) {
	output := make([]any, 0)

	// first find all the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, backend := range allPools {
			if strings.EqualFold(v.(string), *backend.Id) {
				orderedBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
				if err == nil {
					output = append(output, orderedBackendPool)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, backend := range allPools {
		found := false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *backend.Id) {
				found = true
				break
			}
		}

		if !found {
			newBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, newBackendPool)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

type flattenedBackendPoolSettings struct {
	enforceBackendPoolsCertificateNameCheck bool
	backendPoolsSendReceiveTimeoutSeconds   int
}

func combineFrontEndEndpoints(allEndpoints []frontdoors.FrontendEndpoint, orderedIds []any, frontDoorId frontdoors.FrontDoorId) ([]any, error) {
	output := make([]any, 0)

	// first find all the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, frontendEndpoint := range allEndpoints {
			if strings.EqualFold(v.(string), *frontendEndpoint.Id) {
				orderedFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
				if err == nil {
					output = append(output, orderedFrontendEndpoint)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, frontendEndpoint := range allEndpoints {
		found := false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *frontendEndpoint.Id) {
				found = true
				break
			}
		}

		if !found {
			newFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
			if err == nil {
				output = append(output, newFrontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func combineHealthProbeSettingsModel(allHealthProbeSettings []frontdoors.HealthProbeSettingsModel, orderedIds []any, frontDoorId frontdoors.FrontDoorId) []any {
	output := make([]any, 0)

	// first find all the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, healthProbeSetting := range allHealthProbeSettings {
			if strings.EqualFold(v.(string), *healthProbeSetting.Id) {
				orderedHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
				output = append(output, orderedHealthProbeSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, healthProbeSetting := range allHealthProbeSettings {
		found := false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *healthProbeSetting.Id) {
				found = true
				break
			}
		}

		if !found {
			newHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
			output = append(output, newHealthProbeSetting)
		}
	}

	return output
}

func combineLoadBalancingSettingsModel(allLoadBalancingSettings []frontdoors.LoadBalancingSettingsModel, orderedIds []any, frontDoorId frontdoors.FrontDoorId) []any {
	output := make([]any, 0)

	// first find all the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, loadBalancingSetting := range allLoadBalancingSettings {
			if strings.EqualFold(v.(string), *loadBalancingSetting.Id) {
				orderedLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalancingSetting, frontDoorId)
				output = append(output, orderedLoadBalanceSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, loadBalanceSetting := range allLoadBalancingSettings {
		found := false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *loadBalanceSetting.Id) {
				found = true
				break
			}
		}

		if !found {
			newLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalanceSetting, frontDoorId)
			output = append(output, newLoadBalanceSetting)
		}
	}

	return output
}

func combineRoutingRules(allRoutingRules []frontdoors.RoutingRule, oldBlocks any, orderedIds []any, frontDoorId frontdoors.FrontDoorId) ([]any, error) {
	output := make([]any, 0)

	// first find all the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, routingRule := range allRoutingRules {
			if strings.EqualFold(v.(string), *routingRule.Id) {
				orderedRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
				if err == nil {
					output = append(output, orderedRoutingRule)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, routingRule := range allRoutingRules {
		found := false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *routingRule.Id) {
				found = true
				break
			}
		}

		if !found {
			newRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, newRoutingRule)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}
