// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
)

func expandFrontDoorBackendPools(input []any, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.BackendPool {
	if len(input) == 0 {
		return &[]frontdoors.BackendPool{}
	}

	output := make([]frontdoors.BackendPool, 0)

	for _, bp := range input {
		backendPool := bp.(map[string]any)
		backendPoolName := backendPool["name"].(string)
		backendPoolLoadBalancingName := backendPool["load_balancing_name"].(string)
		backendPoolHealthProbeName := backendPool["health_probe_name"].(string)
		backends := backendPool["backend"].([]any)

		backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolName).ID()
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolHealthProbeName).ID()
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolLoadBalancingName).ID()

		result := frontdoors.BackendPool{
			Id:   pointer.To(backendPoolId),
			Name: pointer.To(backendPoolName),
			Properties: &frontdoors.BackendPoolProperties{
				Backends: expandFrontDoorBackend(backends),
				HealthProbeSettings: &frontdoors.SubResource{
					Id: pointer.To(healthProbeId),
				},
				LoadBalancingSettings: &frontdoors.SubResource{
					Id: pointer.To(loadBalancingId),
				},
			},
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorBackend(input []any) *[]frontdoors.Backend {
	if len(input) == 0 {
		return &[]frontdoors.Backend{}
	}

	output := make([]frontdoors.Backend, 0)

	for _, be := range input {
		backend := be.(map[string]any)
		address := backend["address"].(string)
		hostHeader := backend["host_header"].(string)
		enabled := expandFrontDoorBackendEnabledState(backend["enabled"].(bool))
		httpPort := int64(backend["http_port"].(int))
		httpsPort := int64(backend["https_port"].(int))
		priority := int64(backend["priority"].(int))
		weight := int64(backend["weight"].(int))

		result := frontdoors.Backend{
			Address:           pointer.To(address),
			BackendHostHeader: pointer.To(hostHeader),
			EnabledState:      &enabled,
			HTTPPort:          pointer.To(httpPort),
			HTTPSPort:         pointer.To(httpsPort),
			Priority:          pointer.To(priority),
			Weight:            pointer.To(weight),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorBackendEnabledState(isEnabled bool) frontdoors.BackendEnabledState {
	if isEnabled {
		return frontdoors.BackendEnabledStateEnabled
	}
	return frontdoors.BackendEnabledStateDisabled
}

func expandFrontDoorBackendPoolsSettings(enforceCertificateNameCheck bool, backendPoolsSendReceiveTimeoutSeconds int64) *frontdoors.BackendPoolsSettings {
	enforceCheck := frontdoors.EnforceCertificateNameCheckEnabledStateDisabled

	if enforceCertificateNameCheck {
		enforceCheck = frontdoors.EnforceCertificateNameCheckEnabledStateEnabled
	}

	result := frontdoors.BackendPoolsSettings{
		EnforceCertificateNameCheck: &enforceCheck,
		SendRecvTimeoutSeconds:      pointer.To(backendPoolsSendReceiveTimeoutSeconds),
	}

	return &result
}

func expandFrontDoorFrontendEndpoint(input []any, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.FrontendEndpoint {
	if len(input) == 0 {
		return &[]frontdoors.FrontendEndpoint{}
	}

	output := make([]frontdoors.FrontendEndpoint, 0)

	for _, frontendEndpoints := range input {
		frontendEndpoint := frontendEndpoints.(map[string]any)
		hostName := frontendEndpoint["host_name"].(string)
		isSessionAffinityEnabled := frontendEndpoint["session_affinity_enabled"].(bool)
		sessionAffinityTtlSeconds := int64(frontendEndpoint["session_affinity_ttl_seconds"].(int))
		waf := frontendEndpoint["web_application_firewall_policy_link_id"].(string)
		name := frontendEndpoint["name"].(string)
		id := frontdoors.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
		sessionAffinityEnabled := frontdoors.SessionAffinityEnabledStateDisabled

		if isSessionAffinityEnabled {
			sessionAffinityEnabled = frontdoors.SessionAffinityEnabledStateEnabled
		}

		result := frontdoors.FrontendEndpoint{
			Id:   pointer.To(id),
			Name: pointer.To(name),
			Properties: &frontdoors.FrontendEndpointProperties{
				HostName:                    pointer.To(hostName),
				SessionAffinityEnabledState: &sessionAffinityEnabled,
				SessionAffinityTtlSeconds:   pointer.To(sessionAffinityTtlSeconds),
			},
		}

		if waf != "" {
			result.Properties.WebApplicationFirewallPolicyLink = &frontdoors.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink{
				Id: pointer.To(waf),
			}
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorHealthProbeSettingsModel(input []any, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.HealthProbeSettingsModel {
	if len(input) == 0 {
		return &[]frontdoors.HealthProbeSettingsModel{}
	}

	output := make([]frontdoors.HealthProbeSettingsModel, 0)

	for _, hps := range input {
		v := hps.(map[string]any)
		path := v["path"].(string)
		intervalInSeconds := int64(v["interval_in_seconds"].(int))
		name := v["name"].(string)
		enabled := v["enabled"].(bool)

		healthProbeEnabled := frontdoors.HealthProbeEnabledEnabled
		if !enabled {
			healthProbeEnabled = frontdoors.HealthProbeEnabledDisabled
		}
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		result := frontdoors.HealthProbeSettingsModel{
			Id:   pointer.To(healthProbeId),
			Name: pointer.To(name),
			Properties: &frontdoors.HealthProbeSettingsProperties{
				IntervalInSeconds: pointer.To(intervalInSeconds),
				Path:              pointer.To(path),
				Protocol:          pointer.ToEnum[frontdoors.FrontDoorProtocol](v["protocol"].(string)),
				HealthProbeMethod: pointer.ToEnum[frontdoors.FrontDoorHealthProbeMethod](v["probe_method"].(string)),
				EnabledState:      &healthProbeEnabled,
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandFrontDoorLoadBalancingSettingsModel(input []any, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.LoadBalancingSettingsModel {
	if len(input) == 0 {
		return &[]frontdoors.LoadBalancingSettingsModel{}
	}

	output := make([]frontdoors.LoadBalancingSettingsModel, 0)

	for _, lbs := range input {
		loadBalanceSetting := lbs.(map[string]any)
		name := loadBalanceSetting["name"].(string)
		sampleSize := int64(loadBalanceSetting["sample_size"].(int))
		successfulSamplesRequired := int64(loadBalanceSetting["successful_samples_required"].(int))
		additionalLatencyMilliseconds := int64(loadBalanceSetting["additional_latency_milliseconds"].(int))
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		result := frontdoors.LoadBalancingSettingsModel{
			Id:   pointer.To(loadBalancingId),
			Name: pointer.To(name),
			Properties: &frontdoors.LoadBalancingSettingsProperties{
				SampleSize:                    pointer.To(sampleSize),
				SuccessfulSamplesRequired:     pointer.To(successfulSamplesRequired),
				AdditionalLatencyMilliseconds: pointer.To(additionalLatencyMilliseconds),
			},
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorRoutingRule(input []any, frontDoorId frontdoors.FrontDoorId, rulesEngines *map[string]*frontdoors.SubResource) *[]frontdoors.RoutingRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoors.RoutingRule, 0)

	for _, rr := range input {
		routingRule := rr.(map[string]any)
		name := routingRule["name"].(string)
		frontendEndpoints := routingRule["frontend_endpoints"].([]any)
		acceptedProtocols := routingRule["accepted_protocols"].([]any)
		ptm := routingRule["patterns_to_match"].([]any)
		enabled := frontdoors.RoutingRuleEnabledState(expandFrontDoorEnabledState(routingRule["enabled"].(bool)))

		patternsToMatch := make([]string, 0)
		for _, p := range ptm {
			patternsToMatch = append(patternsToMatch, p.(string))
		}

		var routingConfiguration frontdoors.RouteConfiguration
		if rc := routingRule["redirect_configuration"].([]any); len(rc) != 0 {
			routingConfiguration = expandFrontDoorRedirectConfiguration(rc)
		} else if fc := routingRule["forwarding_configuration"].([]any); len(fc) != 0 {
			routingConfiguration = expandFrontDoorForwardingConfiguration(fc, frontDoorId)
		}
		routingRuleId := parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		currentRoutingRule := frontdoors.RoutingRule{
			Id:   pointer.To(routingRuleId),
			Name: pointer.To(name),
			Properties: &frontdoors.RoutingRuleProperties{
				FrontendEndpoints:  expandFrontDoorFrontEndEndpoints(frontendEndpoints, frontDoorId),
				AcceptedProtocols:  expandFrontDoorAcceptedProtocols(acceptedProtocols),
				PatternsToMatch:    &patternsToMatch,
				EnabledState:       &enabled,
				RouteConfiguration: routingConfiguration,
			},
		}

		// Preserve existing rules engine for this routing rule
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/7455#issuecomment-882769364
		if rulesEngines != nil {
			if rulesEngine, ok := (*rulesEngines)[name]; ok {
				currentRoutingRule.Properties.RulesEngine = rulesEngine
			}
		}

		output = append(output, currentRoutingRule)
	}

	return &output
}

func expandFrontDoorAcceptedProtocols(input []any) *[]frontdoors.FrontDoorProtocol {
	if len(input) == 0 {
		return &[]frontdoors.FrontDoorProtocol{}
	}

	output := make([]frontdoors.FrontDoorProtocol, 0)

	for _, ap := range input {
		result := frontdoors.FrontDoorProtocolHTTPS
		if ap.(string) == string(frontdoors.FrontDoorProtocolHTTP) {
			result = frontdoors.FrontDoorProtocolHTTP
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorFrontEndEndpoints(input []any, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.SubResource {
	if len(input) == 0 {
		return &[]frontdoors.SubResource{}
	}

	output := make([]frontdoors.SubResource, 0)

	for _, name := range input {
		frontendEndpointId := frontdoors.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name.(string)).ID()
		result := frontdoors.SubResource{
			Id: pointer.To(frontendEndpointId),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorEnabledState(enabled bool) frontdoors.FrontDoorEnabledState {
	if enabled {
		return frontdoors.FrontDoorEnabledStateEnabled
	}
	return frontdoors.FrontDoorEnabledStateDisabled
}

func expandFrontDoorRedirectConfiguration(input []any) frontdoors.RedirectConfiguration {
	if len(input) == 0 {
		return frontdoors.RedirectConfiguration{}
	}

	v := input[0].(map[string]any)
	customHost := v["custom_host"].(string)
	customPath := v["custom_path"].(string)
	customFragment := v["custom_fragment"].(string)
	customQueryString := v["custom_query_string"].(string)

	redirectConfiguration := frontdoors.RedirectConfiguration{
		CustomHost:       pointer.To(customHost),
		RedirectType:     pointer.ToEnum[frontdoors.FrontDoorRedirectType](v["redirect_type"].(string)),
		RedirectProtocol: pointer.ToEnum[frontdoors.FrontDoorRedirectProtocol](v["redirect_protocol"].(string)),
	}
	// The way the API works is if you don't include the attribute in the structure
	// it is treated as Preserve instead of Replace...
	if customHost != "" {
		redirectConfiguration.CustomHost = pointer.To(customHost)
	}
	if customPath != "" {
		redirectConfiguration.CustomPath = pointer.To(customPath)
	}
	if customFragment != "" {
		redirectConfiguration.CustomFragment = pointer.To(customFragment)
	}
	if customQueryString != "" {
		redirectConfiguration.CustomQueryString = pointer.To(customQueryString)
	}
	return redirectConfiguration
}

func expandFrontDoorForwardingConfiguration(input []any, frontDoorId frontdoors.FrontDoorId) frontdoors.ForwardingConfiguration {
	if len(input) == 0 {
		return frontdoors.ForwardingConfiguration{}
	}

	v := input[0].(map[string]any)
	customForwardingPath := v["custom_forwarding_path"].(string)
	backendPoolName := v["backend_pool_name"].(string)
	cacheUseDynamicCompression := v["cache_use_dynamic_compression"].(bool)
	cacheQueryParameterStripDirective := frontdoors.FrontDoorQuery(v["cache_query_parameter_strip_directive"].(string))
	cacheQueryParameters := v["cache_query_parameters"].([]any)
	cacheDuration := v["cache_duration"].(string)
	cacheEnabled := v["cache_enabled"].(bool)

	// convert list of cache_query_parameters into an array into a comma-separated list
	queryParametersArray := make([]string, 0)
	for _, p := range cacheQueryParameters {
		queryParametersArray = append(queryParametersArray, p.(string))
	}
	queryParametersString := strings.Join(queryParametersArray, ",")

	backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolName).ID()
	backend := &frontdoors.SubResource{
		Id: pointer.To(backendPoolId),
	}

	forwardingConfiguration := frontdoors.ForwardingConfiguration{
		ForwardingProtocol: pointer.ToEnum[frontdoors.FrontDoorForwardingProtocol](v["forwarding_protocol"].(string)),
		BackendPool:        backend,
	}
	// Per the portal, if you enable the cache the cache_query_parameter_strip_directive
	// is then a required attribute else the CacheConfiguration type is null
	if cacheEnabled {
		// Set the default value for dynamic compression or use the value defined in the config
		dynamicCompression := frontdoors.DynamicCompressionEnabledEnabled
		if !cacheUseDynamicCompression {
			dynamicCompression = frontdoors.DynamicCompressionEnabledDisabled
		}
		if cacheQueryParameterStripDirective == "" {
			// Set Default Value for strip directive is not in the key slice and cache is enabled
			cacheQueryParameterStripDirective = frontdoors.FrontDoorQueryStripAll
		}
		// set cacheQueryParameters to "" when StripDirective is "StripAll" or "StripNone"
		if cacheQueryParameterStripDirective == "StripAll" || cacheQueryParameterStripDirective == "StripNone" {
			queryParametersString = ""
		}

		// Making sure that duration is empty when cacheDuration is empty
		var duration *string
		if cacheDuration == "" {
			duration = nil
		} else {
			duration = pointer.To(cacheDuration)
		}

		forwardingConfiguration.CacheConfiguration = &frontdoors.CacheConfiguration{
			DynamicCompression:           &dynamicCompression,
			QueryParameterStripDirective: &cacheQueryParameterStripDirective,
			QueryParameters:              pointer.To(queryParametersString),
			CacheDuration:                duration,
		}
	}

	if customForwardingPath != "" {
		forwardingConfiguration.CustomForwardingPath = pointer.To(customForwardingPath)
	}

	return forwardingConfiguration
}
