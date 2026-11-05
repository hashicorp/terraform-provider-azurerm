// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model cognitiveDeploymentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.DeploymentsClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := deployments.NewDeploymentID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &deployments.Deployment{
				Properties: &deployments.DeploymentProperties{},
			}

			properties.Properties.Model = expandDeploymentModelModel(model.Model)

			if model.RaiPolicyName != "" {
				properties.Properties.RaiPolicyName = &model.RaiPolicyName
			}

			if model.DynamicThrottlingEnabled {
				properties.Properties.DynamicThrottlingEnabled = &model.DynamicThrottlingEnabled
			}

			if model.VersionUpgradeOption != "" {
				properties.Properties.VersionUpgradeOption = pointer.ToEnum[deployments.DeploymentModelVersionUpgradeOption](model.VersionUpgradeOption)
			}

			if model.SpilloverDeploymentName != "" {
				properties.Properties.SpilloverDeploymentName = pointer.To(model.SpilloverDeploymentName)
			}

			properties.Sku = expandDeploymentSkuModel(model.Sku)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandDeploymentModelModel(inputList []DeploymentModelModel) *deployments.DeploymentModel {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := deployments.DeploymentModel{}

	if input.Format != "" {
		output.Format = &input.Format
	}

	if input.Name != "" {
		output.Name = &input.Name
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func expandDeploymentSkuModel(inputList []DeploymentSkuModel) *deployments.Sku {
	if len(inputList) == 0 {
		return nil
	}
	input := inputList[0]
	s := &deployments.Sku{
		Name: input.Name,
	}
	if input.Capacity != 0 {
		s.Capacity = pointer.To(input.Capacity)
	}
	if input.Family != "" {
		s.Family = pointer.To(input.Family)
	}
	if input.Size != "" {
		s.Size = pointer.To(input.Size)
	}
	if input.Tier != "" {
		s.Tier = pointer.ToEnum[deployments.SkuTier](input.Tier)
	}
	return s
}
