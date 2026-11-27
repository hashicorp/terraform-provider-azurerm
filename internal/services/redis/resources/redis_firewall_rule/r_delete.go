// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	redisClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisFirewallRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := redisClient.From(meta).FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := redisfirewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.FirewallRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
