// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmTrafficManagerProfileFlatten(d *pluginsdk.ResourceData, id *profiles.TrafficManagerProfileId, model *profiles.Profile) error {
	d.Set("name", id.TrafficManagerProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		if profile := model.Properties; profile != nil {
			profileStatus := ""
			if profile.ProfileStatus != nil {
				profileStatus = string(*profile.ProfileStatus)
			}
			d.Set("profile_status", profileStatus)
			trafficRoutingMethod := ""
			if profile.TrafficRoutingMethod != nil {
				trafficRoutingMethod = string(*profile.TrafficRoutingMethod)
			}
			d.Set("traffic_routing_method", trafficRoutingMethod)
			d.Set("max_return", profile.MaxReturn)

			d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DnsConfig))
			d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))
			trafficViewEnabled := false
			if profile.TrafficViewEnrollmentStatus != nil {
				trafficViewEnabled = *profile.TrafficViewEnrollmentStatus == profiles.TrafficViewEnrollmentStatusEnabled
			}
			d.Set("traffic_view_enabled", trafficViewEnabled)

			// fqdn is actually inside DNSConfig, inlined for simpler reference
			if dns := profile.DnsConfig; dns != nil {
				d.Set("fqdn", dns.Fqdn)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenArmTrafficManagerCustomHeadersConfig(input *[]profiles.MonitorConfigCustomHeadersItem) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}

	headers := *input
	if len(headers) == 0 {
		return result
	}

	for _, v := range headers {
		header := make(map[string]string, 2)
		header["name"] = *v.Name
		header["value"] = *v.Value
		result = append(result, header)
	}

	return result
}

func flattenAzureRMTrafficManagerProfileDNSConfig(dns *profiles.DnsConfig) []any {
	result := make(map[string]any)

	result["relative_name"] = *dns.RelativeName
	result["ttl"] = int(*dns.Ttl)

	return []any{result}
}

func flattenAzureRMTrafficManagerProfileMonitorConfig(cfg *profiles.MonitorConfig) []any {
	result := make(map[string]any)

	protocol := ""
	if cfg.Protocol != nil {
		protocol = string(*cfg.Protocol)
	}
	result["protocol"] = protocol
	result["port"] = int(*cfg.Port)
	result["custom_header"] = flattenArmTrafficManagerCustomHeadersConfig(cfg.CustomHeaders)

	if cfg.Path != nil {
		result["path"] = *cfg.Path
	}

	result["interval_in_seconds"] = int(*cfg.IntervalInSeconds)
	result["timeout_in_seconds"] = int(*cfg.TimeoutInSeconds)
	result["tolerated_number_of_failures"] = int(*cfg.ToleratedNumberOfFailures)

	if v := cfg.ExpectedStatusCodeRanges; v != nil {
		ranges := make([]string, 0)
		for _, r := range *v {
			if r.Min == nil || r.Max == nil {
				continue
			}

			ranges = append(ranges, fmt.Sprintf("%d-%d", *r.Min, *r.Max))
		}
		result["expected_status_code_ranges"] = ranges
	}

	return []any{result}
}
