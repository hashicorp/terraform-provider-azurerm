// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := networkmanagers.ParseNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.NetworkManager.NetworkManagers
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			var state ManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Model.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("scope") {
				existing.Model.Properties.NetworkManagerScopes = expandNetworkManagerScope(state.Scope)
			}

			if metadata.ResourceData.HasChange("scope_accesses") {
				existing.Model.Properties.NetworkManagerScopeAccesses = expandNetworkManagerScopeAccesses(state.ScopeAccesses)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pluginsdk.ExpandPtrMapStringString(state.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
