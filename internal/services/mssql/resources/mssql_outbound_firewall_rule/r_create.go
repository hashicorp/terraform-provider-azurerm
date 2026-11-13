// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_outbound_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlOutboundFirewallRuleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	serverId, err := commonids.ParseSqlServerIDInsensitively(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID %q: %+v", d.Get("server_id"), err)
	}

	id := outboundfirewallrules.NewOutboundFirewallRuleID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_mssql_outbound_firewall_rule", id.ID())
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceMsSqlOutboundFirewallRuleRead(d, meta)
}
