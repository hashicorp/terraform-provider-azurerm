// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_firewall_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceRedisFirewallRuleFlatten(d *pluginsdk.ResourceData, id *redisfirewallrules.FirewallRuleId, firewallRule *redisfirewallrules.RedisFirewallRule) error {
	d.Set("name", id.FirewallRuleName)
	d.Set("redis_cache_name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if firewallRule != nil {
		d.Set("end_ip", firewallRule.Properties.EndIP)
		d.Set("start_ip", firewallRule.Properties.StartIP)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
