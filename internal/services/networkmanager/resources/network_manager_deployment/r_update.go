// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			client := metadata.Client.NetworkManager.NetworkManagers

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{id.Location},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam, networkmanagers.DefaultNetworkManagerDeploymentStatusListOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if resp.Model.Value == nil || len(*resp.Model.Value) == 0 || *(*resp.Model.Value)[0].ConfigurationIds == nil || len(*(*resp.Model.Value)[0].ConfigurationIds) == 0 {
				return metadata.MarkAsGone(id)
			}

			if len(*resp.Model.Value) > 1 {
				return fmt.Errorf("found more than one deployment with id %s", *id)
			}

			deployment := (*resp.Model.Value)[0]
			if deployment.ConfigurationIds == nil {
				return fmt.Errorf("unexpected null configuration ID of %s", *id)
			}

			var state ManagerDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("configuration_ids") {
				deployment.ConfigurationIds = &state.ConfigurationIds
			}

			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: deployment.ConfigurationIds,
				TargetLocations:  []string{state.Location},
				CommitType:       networkmanagers.ConfigurationType(state.ScopeAccess),
			}

			if _, err := client.NetworkManagerCommitsPost(ctx, networkManagerId, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			if err = resourceManagerDeploymentWaitForFinished(ctx, client, id, time.Until(deadline)); err != nil {
				return err
			}

			return nil
		},
		Timeout: 24 * time.Hour,
	}
}
