// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerFirewallRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Postgres.FlexibleServerFirewallRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
