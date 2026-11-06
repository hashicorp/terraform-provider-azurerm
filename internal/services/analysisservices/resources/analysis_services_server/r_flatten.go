// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAnalysisServicesServerFlatten(d *pluginsdk.ResourceData, id *servers.ServerId, model *servers.AnalysisServicesServer) error {
	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			adminUsers := make([]string, 0)
			if props.AsAdministrators != nil && props.AsAdministrators.Members != nil {
				adminUsers = *props.AsAdministrators.Members
			}
			d.Set("admin_users", adminUsers)

			enablePowerBi, fwRules := flattenAnalysisServicesServerFirewallSettings(props)
			d.Set("power_bi_service_enabled", enablePowerBi)
			if err := d.Set("ipv4_firewall_rule", fwRules); err != nil {
				return fmt.Errorf("setting `ipv4_firewall_rule`: %s", err)
			}

			connectionMode := ""
			if props.QuerypoolConnectionMode != nil {
				connectionMode = string(*props.QuerypoolConnectionMode)
			}
			d.Set("querypool_connection_mode", connectionMode)
			d.Set("server_full_name", props.ServerFullName)

			if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
				d.Set("backup_blob_container_uri", containerUri)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAnalysisServicesServerFirewallSettings(serverProperties *servers.AnalysisServicesServerProperties) (*bool, *pluginsdk.Set) {
	if serverProperties == nil || serverProperties.IPV4FirewallSettings == nil {
		return pointer.To(false), pluginsdk.NewSet(hashAnalysisServicesServerIPv4FirewallRule, make([]any, 0))
	}

	firewallSettings := serverProperties.IPV4FirewallSettings

	enablePowerBi := pointer.To(false)
	if firewallSettings.EnablePowerBIService != nil {
		enablePowerBi = firewallSettings.EnablePowerBIService
	}

	fwRules := make([]any, 0)
	if firewallSettings.FirewallRules != nil {
		for _, fwRule := range *firewallSettings.FirewallRules {
			output := make(map[string]any)
			if fwRule.FirewallRuleName != nil {
				output["name"] = *fwRule.FirewallRuleName
			}

			if fwRule.RangeStart != nil {
				output["range_start"] = *fwRule.RangeStart
			}

			if fwRule.RangeEnd != nil {
				output["range_end"] = *fwRule.RangeEnd
			}

			fwRules = append(fwRules, output)
		}
	}

	return enablePowerBi, pluginsdk.NewSet(hashAnalysisServicesServerIPv4FirewallRule, fwRules)
}
