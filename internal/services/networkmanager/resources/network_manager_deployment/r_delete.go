// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagers
			id, err := parse.NetworkManagerDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			input := networkmanagers.NetworkManagerCommit{
				ConfigurationIds: &[]string{},
				TargetLocations:  []string{id.Location},
				CommitType:       networkmanagers.ConfigurationType(id.ScopeAccess),
			}

			networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

			if _, err := client.NetworkManagerCommitsPost(ctx, networkManagerId, input); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			statusClient := metadata.Client.NetworkManager.NetworkManagers
			if err = resourceManagerDeploymentWaitForDeleted(ctx, statusClient, id, time.Until(deadline)); err != nil {
				return err
			}

			return nil
		},
		Timeout: 24 * time.Hour,
	}
}
