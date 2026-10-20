// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AIFoundryProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			hubId, err := workspaces.ParseWorkspaceID(model.AIServicesHubId)
			if err != nil {
				return err
			}

			id := workspaces.NewWorkspaceID(subscriptionId, hubId.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_ai_foundry_project", id.ID())
				}
			}

			payload := workspaces.Workspace{
				Name:     pointer.To(id.WorkspaceName),
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     tags.Expand(model.Tags),
				Kind:     pointer.To("Project"),
				Properties: &workspaces.WorkspaceProperties{
					HubResourceId: pointer.To(hubId.ID()),
				},
			}

			if len(model.Identity) > 0 {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if model.PrimaryUserAssignedIdentity != "" {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(model.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if model.FriendlyName != "" {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}

			if model.HighBusinessImpactEnabled {
				payload.Properties.HbiWorkspace = pointer.To(model.HighBusinessImpactEnabled)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
