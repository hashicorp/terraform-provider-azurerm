// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_firewall_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePostgresqlFlexibleServerFirewallRuleFlatten(d *pluginsdk.ResourceData, id *firewallrules.FirewallRuleId, model *firewallrules.FirewallRule) error {
	d.Set("name", id.FirewallRuleName)
	d.Set("server_id", firewallrules.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID())

	if model != nil {
		d.Set("end_ip_address", model.Properties.EndIPAddress)
		d.Set("start_ip_address", model.Properties.StartIPAddress)
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
