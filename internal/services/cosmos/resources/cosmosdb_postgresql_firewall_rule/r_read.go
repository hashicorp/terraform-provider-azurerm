// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_firewall_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbPostgreSQLFirewallRuleResourceModel{
				Name:      id.FirewallRuleName,
				ClusterId: firewallrules.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
			}

			properties := &model.Properties
			state.EndIPAddress = properties.EndIPAddress
			state.StartIPAddress = properties.StartIPAddress

			return metadata.Encode(&state)
		},
	}
}
