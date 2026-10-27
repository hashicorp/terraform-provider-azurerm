// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_network_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/virtualnetworkrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkrules.ParseVirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
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

	d.Set("name", id.VirtualNetworkRuleName)

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)

			subnetId, err := commonids.ParseSubnetIDInsensitively(props.VirtualNetworkSubnetId)
			if err != nil {
				return fmt.Errorf("parsing subnet ID returned by API %q: %+v", props.VirtualNetworkSubnetId, err)
			}

			d.Set("subnet_id", subnetId.ID())
		}
	}

	return nil
}
