// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_firewall_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMySqlFlexibleServerFirewallRuleFlatten(d *pluginsdk.ResourceData, id *firewallrules.FirewallRuleId, rule *firewallrules.FirewallRule) error {
	d.Set("name", id.FirewallRuleName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if rule != nil {
		d.Set("start_ip_address", rule.Properties.StartIPAddress)
		d.Set("end_ip_address", rule.Properties.EndIPAddress)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
