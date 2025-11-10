// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MongoClusterFirewallRuleResource struct{}

var _ sdk.ResourceWithUpdate = MongoClusterFirewallRuleResource{}

type MongoClusterFirewallRuleResourceModel struct {
	MongoClusterId string `tfschema:"mongo_cluster_id"`
	Name           string `tfschema:"name"`
	EndIpAddress   string `tfschema:"end_ip_address"`
	StartIpAddress string `tfschema:"start_ip_address"`
}

func (r MongoClusterFirewallRuleResource) ModelObject() interface{} {
	return &MongoClusterFirewallRuleResourceModel{}
}

func (r MongoClusterFirewallRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (r MongoClusterFirewallRuleResource) ResourceType() string {
	return "azurerm_mongo_cluster_firewall_rule"
}

func (r MongoClusterFirewallRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.\-_]{0,79}$`),
				"`name` must be between 1 and 80 characters. It must start with an alphanumeric character and can contain alphanumeric characters, dots (.), hyphens (-), and underscores (_).",
			),
		},

		"mongo_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&firewallrules.MongoClusterId{}),

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

func (r MongoClusterFirewallRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MongoClusterFirewallRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.FirewallRulesClient

			var state MongoClusterFirewallRuleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			mongoClusterId, err := firewallrules.ParseMongoClusterID(state.MongoClusterId)
			if err != nil {
				return err
			}

			id := firewallrules.NewFirewallRuleID(mongoClusterId.SubscriptionId, mongoClusterId.ResourceGroupName, mongoClusterId.MongoClusterName, state.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := firewallrules.FirewallRule{
				Properties: &firewallrules.FirewallRuleProperties{
					EndIPAddress:   state.EndIpAddress,
					StartIPAddress: state.StartIpAddress,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MongoClusterFirewallRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MongoClusterFirewallRuleResourceModel{
				Name:           id.FirewallRuleName,
				MongoClusterId: firewallrules.NewMongoClusterID(id.SubscriptionId, id.ResourceGroupName, id.MongoClusterName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.EndIpAddress = props.EndIPAddress
					state.StartIpAddress = props.StartIPAddress
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MongoClusterFirewallRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MongoClusterFirewallRuleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := firewallrules.FirewallRule{
				Properties: &firewallrules.FirewallRuleProperties{
					EndIPAddress:   state.EndIpAddress,
					StartIPAddress: state.StartIpAddress,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MongoClusterFirewallRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
