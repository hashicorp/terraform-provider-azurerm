package check

import "strings"

var skipProps = []string{
	"all.timezone",
	"all.time_zone",
	"all.time_zone_id",
	// "all.serialization.field_delimiter",
	// "azurerm_virtual_network_gateway_connection.type",
	"azurerm_nginx_deployment.identity.type", // there is a diff between real supported values and common identity schema
	// "azurerm_managed_disk.create_option",
	// "azurerm_synapse_role_assignment.role_name",
	// "azurerm_orchestrated_virtual_machine_scale_set.os_profile.linux_configuration.secret.certificate.store",
	// "all.advanced_filter",
}

// skip auto-generated resources document
var skipResource = []string{
	"azurerm_load_test",
	"azurerm_kubernetes_fleet_manager",
}

var skipPropMap = map[string]struct{}{}
var skipResourceMap = map[string]struct{}{}

func init() {
	for _, k := range skipProps {
		skipPropMap[k] = struct{}{}
	}
	for _, k := range skipResource {
		skipResourceMap[k] = struct{}{}
	}
}

func isSkipProp(rt, prop string) bool {
	if _, ok := skipPropMap[rt]; ok {
		return true
	}
	if _, ok := skipPropMap[rt+"."+prop]; ok {
		return true
	}
	if idx := strings.LastIndex(prop, "."); idx > 0 {
		prop = prop[idx+1:]
	}
	if _, ok := skipPropMap["all."+prop]; ok {
		return true
	}
	return false
}
