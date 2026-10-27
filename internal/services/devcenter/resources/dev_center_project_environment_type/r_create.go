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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterProjectEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(model.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := environmenttypes.NewEnvironmentTypeID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ProjectEnvironmentTypesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return err
			}

			parameters := environmenttypes.ProjectEnvironmentType{
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &environmenttypes.ProjectEnvironmentTypeProperties{
					DeploymentTargetId: pointer.To(model.DeploymentTargetId),
					CreatorRoleAssignment: &environmenttypes.ProjectEnvironmentTypeUpdatePropertiesCreatorRoleAssignment{
						Roles: expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(model.CreatorRoleAssignmentRoles),
					},
					Status: pointer.To(environmenttypes.EnvironmentTypeEnableStatusEnabled),
				},
				Identity: identity,
				Tags:     pointer.To(model.Tags),
			}

			userRoleAssignment, err := expandDevCenterProjectEnvironmentTypeUserRoleAssignment(model.UserRoleAssignment)
			if err != nil {
				return err
			}
			parameters.Properties.UserRoleAssignments = userRoleAssignment

			if _, err := client.ProjectEnvironmentTypesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(input []string) *map[string]environmenttypes.EnvironmentRole {
	if len(input) == 0 {
		return nil
	}

	result := map[string]environmenttypes.EnvironmentRole{}

	for _, v := range input {
		result[v] = environmenttypes.EnvironmentRole{}
	}

	return &result
}

func expandDevCenterProjectEnvironmentTypeUserRoleAssignment(input []DevCenterProjectEnvironmentTypeUserRoleAssignment) (*map[string]environmenttypes.UserRoleAssignment, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := map[string]environmenttypes.UserRoleAssignment{}

	for _, v := range input {
		if _, exists := result[v.UserId]; exists {
			return nil, fmt.Errorf("`user_id` is duplicate")
		}

		result[v.UserId] = environmenttypes.UserRoleAssignment{
			Roles: expandDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(v.Roles),
		}
	}

	return &result, nil
}

func expandDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(input []string) *map[string]environmenttypes.EnvironmentRole {
	if len(input) == 0 {
		return nil
	}

	result := map[string]environmenttypes.EnvironmentRole{}

	for _, v := range input {
		result[v] = environmenttypes.EnvironmentRole{}
	}

	return &result
}
