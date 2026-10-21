// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandNetworkConnectionMonitorEndpoint(input []any) *[]connectionmonitors.ConnectionMonitorEndpoint {
	results := make([]connectionmonitors.ConnectionMonitorEndpoint, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := connectionmonitors.ConnectionMonitorEndpoint{
			Name:   v["name"].(string),
			Filter: expandNetworkConnectionMonitorEndpointFilter(v["filter"].([]any)),
		}

		if address := v["address"]; address != "" {
			result.Address = pointer.To(address.(string))
		}

		if coverageLevel := v["coverage_level"]; coverageLevel != "" {
			result.CoverageLevel = pointer.ToEnum[connectionmonitors.CoverageLevel](coverageLevel.(string))
		}

		excludedItems := v["excluded_ip_addresses"].(*pluginsdk.Set).List()
		includedItems := v["included_ip_addresses"].(*pluginsdk.Set).List()
		if len(excludedItems) != 0 || len(includedItems) != 0 {
			result.Scope = &connectionmonitors.ConnectionMonitorEndpointScope{}

			if len(excludedItems) != 0 {
				var excludedAddresses []connectionmonitors.ConnectionMonitorEndpointScopeItem
				for _, v := range excludedItems {
					excludedAddresses = append(excludedAddresses, connectionmonitors.ConnectionMonitorEndpointScopeItem{
						Address: pointer.To(v.(string)),
					})
				}
				result.Scope.Exclude = &excludedAddresses
			}

			if len(includedItems) != 0 {
				var includedAddresses []connectionmonitors.ConnectionMonitorEndpointScopeItem
				for _, v := range includedItems {
					includedAddresses = append(includedAddresses, connectionmonitors.ConnectionMonitorEndpointScopeItem{
						Address: pointer.To(v.(string)),
					})
				}
				result.Scope.Include = &includedAddresses
			}
		}

		if resourceId := v["target_resource_id"]; resourceId != "" {
			result.ResourceId = pointer.To(resourceId.(string))
		}

		if endpointType := v["target_resource_type"]; endpointType != "" {
			result.Type = pointer.ToEnum[connectionmonitors.EndpointType](endpointType.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorEndpointFilter(input []any) *connectionmonitors.ConnectionMonitorEndpointFilter {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &connectionmonitors.ConnectionMonitorEndpointFilter{
		Type:  pointer.ToEnum[connectionmonitors.ConnectionMonitorEndpointFilterType](v["type"].(string)),
		Items: expandNetworkConnectionMonitorEndpointFilterItem(v["item"].(*pluginsdk.Set).List()),
	}
}

func expandNetworkConnectionMonitorEndpointFilterItem(input []any) *[]connectionmonitors.ConnectionMonitorEndpointFilterItem {
	if len(input) == 0 {
		return nil
	}

	results := make([]connectionmonitors.ConnectionMonitorEndpointFilterItem, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := connectionmonitors.ConnectionMonitorEndpointFilterItem{
			Type: pointer.ToEnum[connectionmonitors.ConnectionMonitorEndpointFilterItemType](v["type"].(string)),
		}

		if address := v["address"]; address != "" {
			result.Address = pointer.To(address.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorTestConfiguration(input []any) *[]connectionmonitors.ConnectionMonitorTestConfiguration {
	results := make([]connectionmonitors.ConnectionMonitorTestConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := connectionmonitors.ConnectionMonitorTestConfiguration{
			Name:              v["name"].(string),
			HTTPConfiguration: expandNetworkConnectionMonitorHTTPConfiguration(v["http_configuration"].([]any)),
			IcmpConfiguration: expandNetworkConnectionMonitorIcmpConfiguration(v["icmp_configuration"].([]any)),
			Protocol:          connectionmonitors.ConnectionMonitorTestConfigurationProtocol(v["protocol"].(string)),
			SuccessThreshold:  expandNetworkConnectionMonitorSuccessThreshold(v["success_threshold"].([]any)),
			TcpConfiguration:  expandNetworkConnectionMonitorTCPConfiguration(v["tcp_configuration"].([]any)),
			TestFrequencySec:  pointer.To(int64(v["test_frequency_in_seconds"].(int))),
		}

		if preferredIPVersion := v["preferred_ip_version"]; preferredIPVersion != "" {
			result.PreferredIPVersion = pointer.ToEnum[connectionmonitors.PreferredIPVersion](preferredIPVersion.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorHTTPConfiguration(input []any) *connectionmonitors.ConnectionMonitorHTTPConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	props := &connectionmonitors.ConnectionMonitorHTTPConfiguration{
		Method:         pointer.ToEnum[connectionmonitors.HTTPConfigurationMethod](v["method"].(string)),
		PreferHTTPS:    pointer.To(v["prefer_https"].(bool)),
		RequestHeaders: expandNetworkConnectionMonitorHTTPHeader(v["request_header"].(*pluginsdk.Set).List()),
	}

	if path := v["path"]; path != "" {
		props.Path = pointer.To(path.(string))
	}

	if port := v["port"]; port != 0 {
		props.Port = pointer.To(int64(port.(int)))
	}

	if ranges := v["valid_status_code_ranges"].(*pluginsdk.Set).List(); len(ranges) != 0 {
		props.ValidStatusCodeRanges = pluginsdk.ExpandStringSlice(ranges)
	}

	return props
}

func expandNetworkConnectionMonitorTCPConfiguration(input []any) *connectionmonitors.ConnectionMonitorTcpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	result := &connectionmonitors.ConnectionMonitorTcpConfiguration{
		Port:              pointer.To(int64(v["port"].(int))),
		DisableTraceRoute: pointer.To(!v["trace_route_enabled"].(bool)),
	}

	if destinationPortBehavior := v["destination_port_behavior"].(string); destinationPortBehavior != "" {
		result.DestinationPortBehavior = pointer.ToEnum[connectionmonitors.DestinationPortBehavior](destinationPortBehavior)
	}

	return result
}

func expandNetworkConnectionMonitorIcmpConfiguration(input []any) *connectionmonitors.ConnectionMonitorIcmpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &connectionmonitors.ConnectionMonitorIcmpConfiguration{
		DisableTraceRoute: pointer.To(!v["trace_route_enabled"].(bool)),
	}
}

func expandNetworkConnectionMonitorSuccessThreshold(input []any) *connectionmonitors.ConnectionMonitorSuccessThreshold {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	return &connectionmonitors.ConnectionMonitorSuccessThreshold{
		ChecksFailedPercent: pointer.To(int64(v["checks_failed_percent"].(int))),
		RoundTripTimeMs:     pointer.To(v["round_trip_time_ms"].(float64)),
	}
}

func expandNetworkConnectionMonitorHTTPHeader(input []any) *[]connectionmonitors.HTTPHeader {
	if len(input) == 0 {
		return nil
	}

	results := make([]connectionmonitors.HTTPHeader, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := connectionmonitors.HTTPHeader{
			Name:  pointer.To(v["name"].(string)),
			Value: pointer.To(v["value"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorTestGroup(input []any) *[]connectionmonitors.ConnectionMonitorTestGroup {
	results := make([]connectionmonitors.ConnectionMonitorTestGroup, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := connectionmonitors.ConnectionMonitorTestGroup{
			Name:               v["name"].(string),
			Destinations:       *pluginsdk.ExpandStringSlice(v["destination_endpoints"].(*pluginsdk.Set).List()),
			Disable:            pointer.To(!v["enabled"].(bool)),
			Sources:            *pluginsdk.ExpandStringSlice(v["source_endpoints"].(*pluginsdk.Set).List()),
			TestConfigurations: *pluginsdk.ExpandStringSlice(v["test_configuration_names"].(*pluginsdk.Set).List()),
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorOutput(input []any) *[]connectionmonitors.ConnectionMonitorOutput {
	results := make([]connectionmonitors.ConnectionMonitorOutput, 0)

	for _, item := range input {
		result := connectionmonitors.ConnectionMonitorOutput{
			Type: pointer.To(connectionmonitors.OutputTypeWorkspace),
			WorkspaceSettings: &connectionmonitors.ConnectionMonitorWorkspaceSettings{
				WorkspaceResourceId: pointer.To(item.(string)),
			},
		}

		results = append(results, result)
	}

	return &results
}
