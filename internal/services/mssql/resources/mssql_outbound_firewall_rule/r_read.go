// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_outbound_firewall_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlOutboundFirewallRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outboundfirewallrules.ParseOutboundFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.OutboundFirewallRuleName)

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	return nil
}
