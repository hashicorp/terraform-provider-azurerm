// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_environment_type

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterProjectEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ProjectEnvironmentTypesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			payload := resp.Model

			if metadata.ResourceData.HasChange("creator_role_assignment_roles") {
				payload.Properties.CreatorRoleAssignment.Roles = expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(model.CreatorRoleAssignmentRoles)
			}

			if metadata.ResourceData.HasChange("deployment_target_id") {
				payload.Properties.DeploymentTargetId = pointer.To(model.DeploymentTargetId)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return err
				}
				payload.Identity = identity
			}

			if metadata.ResourceData.HasChange("user_role_assignment") {
				userRoleAssignment, err := expandDevCenterProjectEnvironmentTypeUserRoleAssignment(model.UserRoleAssignment)
				if err != nil {
					return err
				}
				payload.Properties.UserRoleAssignments = userRoleAssignment
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.ProjectEnvironmentTypesCreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
