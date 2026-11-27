// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandArmTrafficManagerMonitorConfig(d *pluginsdk.ResourceData) *profiles.MonitorConfig {
	monitorSets := d.Get("monitor_config").([]any)
	monitor := monitorSets[0].(map[string]any)

	customHeaders := expandArmTrafficManagerCustomHeadersConfig(monitor["custom_header"].([]any))

	cfg := profiles.MonitorConfig{
		Protocol:                  pointer.ToEnum[profiles.MonitorProtocol](monitor["protocol"].(string)),
		CustomHeaders:             customHeaders,
		Port:                      pointer.To(int64(monitor["port"].(int))),
		Path:                      pointer.To(monitor["path"].(string)),
		IntervalInSeconds:         pointer.To(int64(monitor["interval_in_seconds"].(int))),
		TimeoutInSeconds:          pointer.To(int64(monitor["timeout_in_seconds"].(int))),
		ToleratedNumberOfFailures: pointer.To(int64(monitor["tolerated_number_of_failures"].(int))),
	}

	if v, ok := monitor["expected_status_code_ranges"].([]any); ok {
		ranges := make([]profiles.MonitorConfigExpectedStatusCodeRangesItem, 0)
		for _, r := range v {
			parts := strings.Split(r.(string), "-")
			min, _ := strconv.Atoi(parts[0])
			max, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, profiles.MonitorConfigExpectedStatusCodeRangesItem{
				Min: pointer.To(int64(min)),
				Max: pointer.To(int64(max)),
			})
		}
		cfg.ExpectedStatusCodeRanges = &ranges
	}

	return &cfg
}

func expandArmTrafficManagerCustomHeadersConfig(d []any) *[]profiles.MonitorConfigCustomHeadersItem {
	if len(d) == 0 || d[0] == nil {
		return nil
	}

	customHeaders := make([]profiles.MonitorConfigCustomHeadersItem, len(d))

	for i, v := range d {
		ch := v.(map[string]any)
		customHeaders[i] = profiles.MonitorConfigCustomHeadersItem{
			Name:  pointer.To(ch["name"].(string)),
			Value: pointer.To(ch["value"].(string)),
		}
	}

	return &customHeaders
}

func expandArmTrafficManagerDNSConfig(d *pluginsdk.ResourceData) *profiles.DnsConfig {
	dnsSets := d.Get("dns_config").([]any)
	dns := dnsSets[0].(map[string]any)

	return &profiles.DnsConfig{
		RelativeName: pointer.To(dns["relative_name"].(string)),
		Ttl:          pointer.To(int64(dns["ttl"].(int))),
	}
}

func expandArmTrafficManagerTrafficView(s bool) *profiles.TrafficViewEnrollmentStatus {
	enabled := profiles.TrafficViewEnrollmentStatusDisabled
	if s {
		enabled = profiles.TrafficViewEnrollmentStatusEnabled
	}
	return &enabled
}

// expandTrafficManagerProfileStatus is the status to send, from `enabled`, or in 5.x from `profile_status` when that is set. Leaving both out means enabled.
func expandTrafficManagerProfileStatus(d *pluginsdk.ResourceData) *profiles.ProfileStatus {
	enabled := d.Get("enabled").(bool)
	if !features.SixPointOh() {
		if !pluginsdk.IsExplicitlyNullInConfig(d, "profile_status") {
			return pointer.ToEnum[profiles.ProfileStatus](d.Get("profile_status").(string))
		}

		if pluginsdk.IsExplicitlyNullInConfig(d, "enabled") {
			enabled = true
		}
	}

	if enabled {
		return pointer.To(profiles.ProfileStatusEnabled)
	}

	return pointer.To(profiles.ProfileStatusDisabled)
}
