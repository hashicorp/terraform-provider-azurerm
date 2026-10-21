// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ManagerDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.NetworkManager.NetworkManagers

			networkManagerId, err := networkmanagers.ParseNetworkManagerID(state.NetworkManagerId)
			if err != nil {
				return err
			}

			normalizedLocation := location.Normalize(state.Location)
			id := parse.NewNetworkManagerDeploymentID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, normalizedLocation, state.ScopeAccess)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
				Regions:         &[]string{normalizedLocation},
				DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(state.ScopeAccess)},
			}
			resp, err := client.NetworkManagerDeploymentStatusList(ctx, *networkManagerId, listParam, networkmanagers.DefaultNetworkManagerDeploymentStatusListOperationOptions())

			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if !response.WasNotFound(resp.HttpResponse) && resp.Model.Value != nil && len(*resp.Model.Value) != 0 && *(*resp.Model.Value)[0].ConfigurationIds != nil && len(*(*resp.Model.Value)[0].ConfigurationIds) != 0 {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: &state.ConfigurationIds,
				TargetLocations:  []string{state.Location},
				CommitType:       networkmanagers.ConfigurationType(state.ScopeAccess),
			}

			// The API has an issue where it may return an async operation URL that is invalid, breaking the default LRO poller
			if _, err := client.NetworkManagerCommitsPost(ctx, *networkManagerId, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if metadata.Client.Features.PersistIDOnCreateBeforePollingForCompletion {
				metadata.SetID(id)
			}

			poller := pollers.NewPoller(custompollers.NewNetworkManagerDeploymentPoller(client, *id), time.Second*10, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 24 * time.Hour,
	}
}
