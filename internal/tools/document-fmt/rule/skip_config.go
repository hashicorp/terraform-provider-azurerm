// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"regexp"
	"strings"
)

var skipProps = []string{
	"azurerm_load_test",
	"azurerm_kubernetes_fleet_manager",
	"azurerm_dev_center",

	"all.timezone",
	"all.time_zone",
	"all.time_zone_id",
	"azurerm_nginx_deployment.identity.type", // there is a diff between real supported values and common identity schema
	"azurerm_kubernetes_cluster.default_node_pool.os_sku",
	"azurerm_kubernetes_cluster_node_pool.os_sku",
	"/azurerm_container_app.template..+_scale_rule.authentication.trigger_parameter",
	"/azurerm_container_app.template..+_probe.header",
	`/azurerm_automanage_configuration.backup.retention_policy.(?:weekly|daily)_schedule.retention_duration.duration_type`,
	"/azurerm_sentinel_metadata.dependency.*",
	"/azurerm_orchestrated_virtual_machine_scale_set.os_profile.(windows|linux)_configuration.secret.certificate",
	"/azurerm_eventgrid_event_subscription.advanced_filter.*",
	"/azurerm_eventgrid_system_topic_event_subscription.advanced_filter.*",
}

var skipConfig = &struct {
	skipPaths  map[string]struct{}
	skipRegexp []*regexp.Regexp
}{
	skipPaths: map[string]struct{}{},
}

func init() {
	for _, k := range skipProps {
		if strings.HasPrefix(k, "/") {
			skipConfig.skipRegexp = append(skipConfig.skipRegexp, regexp.MustCompile(k[1:]))
		} else {
			skipConfig.skipPaths[k] = struct{}{}
		}
	}
}

// isSkipProp checks if a property should be skipped in validation
// rt is the resource type (e.g., "azurerm_storage_account")
// prop is the property path (e.g., "blob_properties.cors_rule")
func isSkipProp(rt, prop string) bool {
	// Check if entire resource is skipped
	if _, ok := skipConfig.skipPaths[rt]; ok {
		return true
	}

	// Check if specific property path is skipped
	if _, ok := skipConfig.skipPaths[rt+"."+prop]; ok {
		return true
	}

	// Check if property name (last part) matches an "all.xxx" pattern
	allKey := prop
	if idx := strings.LastIndex(prop, "."); idx > 0 {
		allKey = prop[idx+1:]
	}
	if _, ok := skipConfig.skipPaths["all."+allKey]; ok {
		return true
	}

	// Check against regex patterns
	for _, reg := range skipConfig.skipRegexp {
		if reg.MatchString(rt + "." + prop) {
			return true
		}
	}

	return false
}
