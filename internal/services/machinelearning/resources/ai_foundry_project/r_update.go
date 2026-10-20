// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AIFoundryProjectModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			// Hubs and Projects share the same API where Projects inherit the KV/Storage/AppInsights/CR/Network/Encryption settings from the Hub
			// When updating a Project the API will error when trying to send the inherited settings that get returned when we retrieve the resource for patching in changes
			// This is a hack to work around this behaviour and design, so that we can continue to support the use of `ignore_changes` on the resource
			payload.Properties.ManagedNetwork = nil
			payload.Properties.KeyVault = nil
			payload.Properties.StorageAccount = nil
			payload.Properties.ContainerRegistry = nil
			payload.Properties.ApplicationInsights = nil
			payload.Properties.Encryption = nil

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("primary_user_assigned_identity") {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(state.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = pointer.To(state.FriendlyName)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(state.Tags)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
