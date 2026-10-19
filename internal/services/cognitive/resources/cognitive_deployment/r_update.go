// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			id, err := deployments.ParseDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			properties := resp.Model

			if metadata.ResourceData.HasChange("dynamic_throttling_enabled") {
				properties.Properties.DynamicThrottlingEnabled = pointer.To(model.DynamicThrottlingEnabled)
			}

			if metadata.ResourceData.HasChange("sku.0.capacity") {
				properties.Sku.Capacity = pointer.To(model.Sku[0].Capacity)
			}

			if metadata.ResourceData.HasChange("rai_policy_name") {
				properties.Properties.RaiPolicyName = pointer.To(model.RaiPolicyName)
			}

			if metadata.ResourceData.HasChange("model.0.version") {
				properties.Properties.Model.Version = pointer.To(model.Model[0].Version)
			}

			if metadata.ResourceData.HasChange("spillover_deployment_name") {
				properties.Properties.SpilloverDeploymentName = pointer.To(model.SpilloverDeploymentName)
			}

			properties.Properties.VersionUpgradeOption = pointer.ToEnum[deployments.DeploymentModelVersionUpgradeOption](model.VersionUpgradeOption)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
