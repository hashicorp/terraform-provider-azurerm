// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_fqdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.SettingsRuleGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_machine_learning_workspace_network_outbound_rule_fqdn", id.ID())
				}
			}

			outboundRule := managednetwork.OutboundRuleBasicResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(managednetwork.RuleTypeFQDN)),
				Properties: managednetwork.FqdnOutboundRule{
					Category:    pointer.To(managednetwork.RuleCategoryUserDefined),
					Destination: &model.Destination,
				},
			}

			if err = client.SettingsRuleCreateOrUpdateCallbackThenPoll(ctx, id, outboundRule, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
