// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			hub := AIFoundryProjectModel{
				Name: id.WorkspaceName,
			}

			if model := resp.Model; model != nil {
				hub.Location = location.NormalizeNilable(model.Location)
				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				hub.Identity = flattenedIdentity

				hub.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if v := pointer.From(props.HubResourceId); v != "" {
						hubId, err := workspaces.ParseWorkspaceID(v)
						if err != nil {
							return err
						}
						hub.AIServicesHubId = hubId.ID()
					}

					hub.Description = pointer.From(props.Description)
					hub.FriendlyName = pointer.From(props.FriendlyName)
					hub.HighBusinessImpactEnabled = pointer.From(props.HbiWorkspace)
					hub.ProjectId = pointer.From(props.WorkspaceId)
					hub.PrimaryUserAssignedIdentity = pointer.From(props.PrimaryUserAssignedIdentity)
				}
			}

			return metadata.Encode(&hub)
		},
	}
}
