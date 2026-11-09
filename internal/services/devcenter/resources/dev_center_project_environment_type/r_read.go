// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_environment_type

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectEnvironmentTypesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterProjectEnvironmentTypeResourceModel{
				Name:               id.EnvironmentTypeName,
				DevCenterProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(pointer.From(model.Location))
				state.Tags = pointer.From(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					state.DeploymentTargetId = pointer.From(props.DeploymentTargetId)
					state.UserRoleAssignment = flattenDevCenterProjectEnvironmentTypeUserRoleAssignment(props.UserRoleAssignments)

					if v := props.CreatorRoleAssignment; v != nil {
						state.CreatorRoleAssignmentRoles = flattenDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(v.Roles)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

type DevCenterProjectEnvironmentTypeUserRoleAssignment struct {
	UserId string   `tfschema:"user_id"`
	Roles  []string `tfschema:"roles"`
}
