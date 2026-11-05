// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient

			id, err := connectedclusters.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesProvisionedClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ConnectedClusterGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			payload := resp.Model

			if metadata.ResourceData.HasChange("azure_active_directory") {
				payload.Properties.AadProfile = expandArcKubernetesClusterAadProfile(model.AzureActiveDirectory)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("arc_agent_desired_version") {
				if desiredVersion := model.ArcAgentDesiredVersion; desiredVersion != "" {
					payload.Properties.ArcAgentProfile.DesiredAgentVersion = pointer.To(desiredVersion)
				} else {
					payload.Properties.ArcAgentProfile.DesiredAgentVersion = nil
				}
			}

			if metadata.ResourceData.HasChange("arc_agent_auto_upgrade_enabled") {
				autoUpgradeOption := connectedclusters.AutoUpgradeOptionsEnabled
				if !model.ArcAgentAutoUpgradeEnabled {
					autoUpgradeOption = connectedclusters.AutoUpgradeOptionsDisabled
				}

				payload.Properties.ArcAgentProfile.AgentAutoUpgrade = pointer.To(autoUpgradeOption)
			}

			if err := client.ConnectedClusterCreateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
