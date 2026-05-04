// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var resourceTypeSupportSubResType = map[string][]string{
	"Microsoft.KeyVault":                {"vault"},
	"Microsoft.Cache":                   {"redisCache"},
	"Microsoft.MachineLearningServices": {"amlworkspace"},
	"Microsoft.Storage":                 {"blob", "table", "queue", "file", "web", "dfs"},
}

type machineLearningWorkspaceOutboundRulePrivateEndpointModel struct {
	Name              string `tfschema:"name"`
	WorkspaceId       string `tfschema:"workspace_id"`
	ServiceResourceId string `tfschema:"service_resource_id"`
	SubresourceTarget string `tfschema:"sub_resource_target"`
	SparkEnabled      bool   `tfschema:"spark_enabled"`
}

type WorkspaceNetworkOutboundRulePrivateEndpoint struct{}

var _ sdk.Resource = WorkspaceNetworkOutboundRulePrivateEndpoint{}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) ResourceType() string {
	return "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint"
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) ModelObject() interface{} {
	return &machineLearningWorkspaceOutboundRulePrivateEndpointModel{}
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managednetwork.ValidateOutboundRuleID
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) Arguments() map[string]*pluginsdk.Schema {
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

		"service_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sub_resource_target": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"vault",
				"amlworkspace",
				"blob",
				"table",
				"queue",
				"file",
				"web",
				"dfs",
				"redisCache",
			}, false),
		},

		"spark_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
	}
	return arguments
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model machineLearningWorkspaceOutboundRulePrivateEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

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
				return tf.ImportAsExistsError("azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", id.ID())
			}

			resId, err := resourceids.ParseAzureResourceID(model.ServiceResourceId)
			if err != nil {
				return err
			}

			supportType := false
			if subTypes, ok := resourceTypeSupportSubResType[resId.Provider]; ok {
				for _, typ := range subTypes {
					if strings.EqualFold(typ, model.SubresourceTarget) {
						supportType = true
						break
					}
				}
			}

			if !supportType {
				return fmt.Errorf(" unsupported resource type: %s. Sub resource type supported by Service Resource ID: %s is %s ",
					model.SubresourceTarget, model.ServiceResourceId,
					strings.Join(resourceTypeSupportSubResType[resId.Provider], ", "))
			}

			outboundRule := managednetwork.OutboundRuleBasicResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(managednetwork.RuleTypePrivateEndpoint)),
				Properties: managednetwork.PrivateEndpointOutboundRule{
					Category: pointer.To(managednetwork.RuleCategoryUserDefined),
					Destination: &managednetwork.PrivateEndpointDestination{
						ServiceResourceId: pointer.To(model.ServiceResourceId),
						SubresourceTarget: pointer.To(model.SubresourceTarget),
						SparkEnabled:      pointer.To(model.SparkEnabled),
					},
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

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) Read() sdk.ResourceFunc {
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

			state := machineLearningWorkspaceOutboundRulePrivateEndpointModel{
				Name:        id.OutboundRuleName,
				WorkspaceId: managednetwork.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if prop, ok := props.(managednetwork.PrivateEndpointOutboundRule); ok && prop.Destination != nil {
						state.SparkEnabled = pointer.From(prop.Destination.SparkEnabled)
						state.SubresourceTarget = pointer.From(prop.Destination.SubresourceTarget)
						state.ServiceResourceId = pointer.From(prop.Destination.ServiceResourceId)
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r WorkspaceNetworkOutboundRulePrivateEndpoint) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.SettingsRuleDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}
