package postgresqlhsc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PostgreSQLHyperScaleFirewallRuleModel struct {
	Name           string `tfschema:"name"`
	ServerGroupId  string `tfschema:"server_group_id"`
	EndIPAddress   string `tfschema:"end_ip_address"`
	StartIPAddress string `tfschema:"start_ip_address"`
}

type PostgreSQLHyperScaleFirewallRuleResource struct{}

var _ sdk.ResourceWithUpdate = PostgreSQLHyperScaleFirewallRuleResource{}

func (r PostgreSQLHyperScaleFirewallRuleResource) ResourceType() string {
	return "azurerm_postgresql_hyperscale_firewall_rule"
}

func (r PostgreSQLHyperScaleFirewallRuleResource) ModelObject() interface{} {
	return &PostgreSQLHyperScaleFirewallRuleModel{}
}

func (r PostgreSQLHyperScaleFirewallRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (r PostgreSQLHyperScaleFirewallRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: firewallrules.ValidateServerGroupsv2ID,
		},

		"end_ip_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"start_ip_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r PostgreSQLHyperScaleFirewallRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgreSQLHyperScaleFirewallRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PostgreSQLHyperScaleFirewallRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.FirewallRulesClient
			serverGroupId, err := firewallrules.ParseServerGroupsv2ID(model.ServerGroupId)
			if err != nil {
				return err
			}

			id := firewallrules.NewFirewallRuleID(serverGroupId.SubscriptionId, serverGroupId.ResourceGroupName, serverGroupId.ServerGroupsv2Name, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &firewallrules.FirewallRule{
				Properties: firewallrules.FirewallRuleProperties{
					EndIPAddress:   model.EndIPAddress,
					StartIPAddress: model.StartIPAddress,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PostgreSQLHyperScaleFirewallRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.FirewallRulesClient

			id, err := firewallrules.ParseFirewallRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PostgreSQLHyperScaleFirewallRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
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

func (r PostgreSQLHyperScaleFirewallRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.FirewallRulesClient

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

			state := PostgreSQLHyperScaleFirewallRuleModel{
				Name:          id.FirewallRuleName,
				ServerGroupId: firewallrules.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
			}

			properties := &model.Properties
			state.EndIPAddress = properties.EndIPAddress
			state.StartIPAddress = properties.StartIPAddress

			return metadata.Encode(&state)
		},
	}
}

func (r PostgreSQLHyperScaleFirewallRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.FirewallRulesClient

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
