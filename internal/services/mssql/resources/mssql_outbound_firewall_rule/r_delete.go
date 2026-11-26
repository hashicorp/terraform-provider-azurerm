// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_outbound_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlOutboundFirewallRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outboundfirewallrules.ParseOutboundFirewallRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	return nil
}
