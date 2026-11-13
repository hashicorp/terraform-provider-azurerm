// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySqlFlexibleServerFirewallRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.FirewallRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
