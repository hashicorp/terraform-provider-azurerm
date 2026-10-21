// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenNetworkConnectionMonitorEndpoint(input *[]connectionmonitors.ConnectionMonitorEndpoint) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var coverageLevel string
		if item.CoverageLevel != nil && string(*item.CoverageLevel) != "" {
			coverageLevel = string(*item.CoverageLevel)
		}

		var endpointType string
		if item.Type != nil && string(*item.Type) != "" {
			endpointType = string(*item.Type)
		}

		v := map[string]any{
			"name":                 item.Name,
			"address":              pointer.From(item.Address),
			"coverage_level":       coverageLevel,
			"target_resource_id":   pointer.From(item.ResourceId),
			"target_resource_type": endpointType,
			"filter":               flattenNetworkConnectionMonitorEndpointFilter(item.Filter),
		}

		if scope := item.Scope; scope != nil {
			if includeScope := scope.Include; includeScope != nil {
				includedAddresses := make([]any, 0)

				for _, includedItem := range *includeScope {
					if includedAddress := includedItem.Address; includedAddress != nil {
						includedAddresses = append(includedAddresses, includedAddress)
					}
				}

				v["included_ip_addresses"] = includedAddresses
			}

			if excludeScope := scope.Exclude; excludeScope != nil {
				excludedAddresses := make([]any, 0)

				for _, excludedItem := range *excludeScope {
					if excludedAddress := excludedItem.Address; excludedAddress != nil {
						excludedAddresses = append(excludedAddresses, excludedAddress)
					}
				}

				v["excluded_ip_addresses"] = excludedAddresses
			}
		}

		results = append(results, v)
	}
	return results
}

func flattenNetworkConnectionMonitorEndpointFilter(input *connectionmonitors.ConnectionMonitorEndpointFilter) []any {
	if input == nil {
		return make([]any, 0)
	}

	var t connectionmonitors.ConnectionMonitorEndpointFilterType
	if input.Type != nil && string(*input.Type) != "" {
		t = *input.Type
	}
	return []any{
		map[string]any{
			"item": flattenNetworkConnectionMonitorEndpointFilterItem(input.Items),
			"type": t,
		},
	}
}

func flattenNetworkConnectionMonitorEndpointFilterItem(input *[]connectionmonitors.ConnectionMonitorEndpointFilterItem) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var t connectionmonitors.ConnectionMonitorEndpointFilterItemType
		if item.Type != nil && string(*item.Type) != "" {
			t = *item.Type
		}

		v := map[string]any{
			"address": pointer.From(item.Address),
			"type":    t,
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorTestConfiguration(input *[]connectionmonitors.ConnectionMonitorTestConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var protocol connectionmonitors.ConnectionMonitorTestConfigurationProtocol
		if item.Protocol != "" {
			protocol = item.Protocol
		}

		var preferredIpVersion connectionmonitors.PreferredIPVersion
		if item.PreferredIPVersion != nil && string(*item.PreferredIPVersion) != "" {
			preferredIpVersion = *item.PreferredIPVersion
		}

		v := map[string]any{
			"name":                      item.Name,
			"protocol":                  protocol,
			"http_configuration":        flattenNetworkConnectionMonitorHTTPConfiguration(item.HTTPConfiguration),
			"icmp_configuration":        flattenNetworkConnectionMonitorIcmpConfiguration(item.IcmpConfiguration),
			"preferred_ip_version":      preferredIpVersion,
			"success_threshold":         flattenNetworkConnectionMonitorSuccessThreshold(item.SuccessThreshold),
			"tcp_configuration":         flattenNetworkConnectionMonitorTCPConfiguration(item.TcpConfiguration),
			"test_frequency_in_seconds": pointer.From(item.TestFrequencySec),
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorHTTPConfiguration(input *connectionmonitors.ConnectionMonitorHTTPConfiguration) []any {
	if input == nil {
		return make([]any, 0)
	}

	var method connectionmonitors.HTTPConfigurationMethod
	if input.Method != nil && string(*input.Method) != "" {
		method = *input.Method
	}

	return []any{
		map[string]any{
			"method":                   method,
			"path":                     pointer.From(input.Path),
			"port":                     pointer.From(input.Port),
			"prefer_https":             pointer.From(input.PreferHTTPS),
			"request_header":           flattenNetworkConnectionMonitorHTTPHeader(input.RequestHeaders),
			"valid_status_code_ranges": pluginsdk.FlattenSlice(input.ValidStatusCodeRanges),
		},
	}
}

func flattenNetworkConnectionMonitorIcmpConfiguration(input *connectionmonitors.ConnectionMonitorIcmpConfiguration) []any {
	if input == nil {
		return make([]any, 0)
	}

	var enableTraceRoute bool
	if input.DisableTraceRoute != nil {
		enableTraceRoute = !*input.DisableTraceRoute
	}

	return []any{
		map[string]any{
			"trace_route_enabled": enableTraceRoute,
		},
	}
}

func flattenNetworkConnectionMonitorSuccessThreshold(input *connectionmonitors.ConnectionMonitorSuccessThreshold) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"checks_failed_percent": pointer.From(input.ChecksFailedPercent),
			"round_trip_time_ms":    pointer.From(input.RoundTripTimeMs),
		},
	}
}

func flattenNetworkConnectionMonitorTCPConfiguration(input *connectionmonitors.ConnectionMonitorTcpConfiguration) []any {
	if input == nil {
		return make([]any, 0)
	}

	var enableTraceRoute bool
	if input.DisableTraceRoute != nil {
		enableTraceRoute = !*input.DisableTraceRoute
	}

	var destinationPortBehavior connectionmonitors.DestinationPortBehavior
	if input.DestinationPortBehavior != nil && string(*input.DestinationPortBehavior) != "" {
		destinationPortBehavior = *input.DestinationPortBehavior
	}

	return []any{
		map[string]any{
			"trace_route_enabled":       enableTraceRoute,
			"port":                      pointer.From(input.Port),
			"destination_port_behavior": string(destinationPortBehavior),
		},
	}
}

func flattenNetworkConnectionMonitorHTTPHeader(input *[]connectionmonitors.HTTPHeader) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := map[string]any{
			"name":  pointer.From(item.Name),
			"value": pointer.From(item.Value),
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorTestGroup(input *[]connectionmonitors.ConnectionMonitorTestGroup) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := map[string]any{
			"name":                     item.Name,
			"destination_endpoints":    item.Destinations,
			"source_endpoints":         item.Sources,
			"test_configuration_names": item.TestConfigurations,
			"enabled":                  !pointer.From(item.Disable),
		}

		results = append(results, v)
	}
	return results
}

func flattenNetworkConnectionMonitorOutput(input *[]connectionmonitors.ConnectionMonitorOutput) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var workspaceResourceId string
		if item.WorkspaceSettings != nil && item.WorkspaceSettings.WorkspaceResourceId != nil {
			workspaceResourceId = *item.WorkspaceSettings.WorkspaceResourceId
		}

		results = append(results, workspaceResourceId)
	}

	return results
}
