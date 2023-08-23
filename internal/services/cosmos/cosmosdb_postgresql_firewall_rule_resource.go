// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDbPostgreSQLFirewallRuleResourceModel struct {
	Name           string `tfschema:"name"`
	ClusterId      string `tfschema:"cluster_id"`
	EndIPAddress   string `tfschema:"end_ip_address"`
	StartIPAddress string `tfschema:"start_ip_address"`
}

type CosmosDbPostgreSQLFirewallRuleResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbPostgreSQLFirewallRuleResource{}

func (r CosmosDbPostgreSQLFirewallRuleResource) ResourceType() string {
	return "azurerm_cosmosdb_postgresql_firewall_rule"
}

func (r CosmosDbPostgreSQLFirewallRuleResource) ModelObject() interface{} {
	return &CosmosDbPostgreSQLFirewallRuleResourceModel{}
}

func (r CosmosDbPostgreSQLFirewallRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallRuleName,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: firewallrules.ValidateServerGroupsv2ID,
		},

		"end_ip_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"start_ip_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
	}
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLFirewallRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.FirewallRulesClient
			clusterId, err := firewallrules.ParseServerGroupsv2ID(model.ClusterId)
			if err != nil {
				return err
			}

			id := firewallrules.NewFirewallRuleID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ServerGroupsv2Name, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := firewallrules.FirewallRule{
				Properties: firewallrules.FirewallRuleProperties{
					EndIPAddress:   model.EndIPAddress,
					StartIPAddress: model.StartIPAddress,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CosmosDbPostgreSQLFirewallRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("end_ip_address") {
				parameters.Properties.EndIPAddress = model.EndIPAddress
			}

			if metadata.ResourceData.HasChange("start_ip_address") {
				parameters.Properties.StartIPAddress = model.StartIPAddress
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Read() sdk.ResourceFunc {
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

func (r CosmosDbPostgreSQLFirewallRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
