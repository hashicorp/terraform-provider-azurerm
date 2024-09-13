// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = WorkspaceNetworkOutboundRuleFqdn{}

type machineLearningWorkspaceOutboundRuleFqdnModel struct {
	Name        string `tfschema:"name"`
	WorkspaceId string `tfschema:"workspace_id"`
	Destination string `tfschema:"destination"`
}

type WorkspaceNetworkOutboundRuleFqdn struct{}

var _ sdk.Resource = WorkspaceNetworkOutboundRuleFqdn{}

func (r WorkspaceNetworkOutboundRuleFqdn) ResourceType() string {
	return "azurerm_machine_learning_workspace_network_outbound_rule_fqdn"
}

func (r WorkspaceNetworkOutboundRuleFqdn) ModelObject() interface{} {
	return &machineLearningWorkspaceOutboundRuleFqdnModel{}
}

func (r WorkspaceNetworkOutboundRuleFqdn) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managednetwork.ValidateOutboundRuleID
}

func (r WorkspaceNetworkOutboundRuleFqdn) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managednetwork.ValidateWorkspaceID,
		},

		"destination": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return arguments
}

func (r WorkspaceNetworkOutboundRuleFqdn) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkspaceNetworkOutboundRuleFqdn) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceOutboundRuleFqdnModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.ManagedNetwork
			subscriptionId := metadata.Client.Account.SubscriptionId

			workspaceId, err := managednetwork.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return err
			}
			id := managednetwork.NewOutboundRuleID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)
			existing, err := client.SettingsRuleGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_workspace_network_outbound_rule_fqdn", id.ID())
			}

			outboundRule := managednetwork.OutboundRuleBasicResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(managednetwork.RuleTypeFQDN)),
				Properties: managednetwork.FqdnOutboundRule{
					Category:    pointer.To(managednetwork.RuleCategoryUserDefined),
					Destination: &model.Destination,
				},
			}

			if err = client.SettingsRuleCreateOrUpdateThenPoll(ctx, id, outboundRule); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkspaceNetworkOutboundRuleFqdn) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceOutboundRuleFqdnModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.ManagedNetwork
			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.SettingsRuleGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("destination") {
				payload.Properties = managednetwork.FqdnOutboundRule{
					Destination: &model.Destination,
				}
			}

			if err = client.SettingsRuleCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WorkspaceNetworkOutboundRuleFqdn) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.SettingsRuleGet(ctx, *id)
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

			state := machineLearningWorkspaceOutboundRuleFqdnModel{
				Name: id.OutboundRuleName,
			}

			if props := model.Properties; props != nil {
				if prop, ok := props.(managednetwork.FqdnOutboundRule); ok {
					if prop.Destination != nil {
						state.Destination = *prop.Destination
					}
				}
			}

			state.WorkspaceId = workspaces.NewWorkspaceID(metadata.Client.Account.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()

			return metadata.Encode(&state)
		},
	}
}

func (r WorkspaceNetworkOutboundRuleFqdn) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.SettingsRuleDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting Machine Learning Workspace FQDN Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			if err = future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Workspace FQDN Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			return nil
		},
	}
}
