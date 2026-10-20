// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
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

			var state AIFoundryModel
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

			if metadata.ResourceData.HasChange("application_insights_id") {
				applicationInsightsId, err := componentsapis.ParseComponentID(state.ApplicationInsightsId)
				if err != nil {
					return err
				}
				payload.Properties.ApplicationInsights = pointer.To(applicationInsightsId.ID())
			}

			if metadata.ResourceData.HasChange("container_registry_id") {
				containerRegistryId, err := registries.ParseRegistryID(state.ContainerRegistryId)
				if err != nil {
					return err
				}
				payload.Properties.ContainerRegistry = pointer.To(containerRegistryId.ID())
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.ToEnum[workspaces.PublicNetworkAccess](state.PublicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(state.Description)
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

			if metadata.ResourceData.HasChange("primary_user_assigned_identity") {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(state.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if metadata.ResourceData.HasChange("managed_network") {
				payload.Properties.ManagedNetwork = expandManagedNetwork(state.ManagedNetwork)
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
