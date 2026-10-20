// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_private_endpoint

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
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.SettingsRuleGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", id.ID())
				}
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

			if err = client.SettingsRuleCreateOrUpdateCallbackThenPoll(ctx, id, outboundRule, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
