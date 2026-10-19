// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/role_definition/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			stateId, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			id := roledefinitions.NewScopedRoleDefinitionID(stateId.Scope, stateId.RoleID)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(stateId)
				}

				return fmt.Errorf("retrieving %s: %+v", stateId, err)
			}

			state := RoleDefinitionModel{
				Scope:            stateId.Scope,
				RoleDefinitionId: stateId.RoleID,
			}

			if model := resp.Model; model != nil {
				// The Azure resource id of Role Definition is not as same as the one we used to create it.
				// So we read from the response.
				state.RoleDefinitionResourceId = pointer.From(model.Id)
				if prop := model.Properties; prop != nil {
					state.Name = pointer.From(model.Properties.RoleName)
					state.Description = pointer.From(prop.Description)
					state.Permissions = flattenRoleDefinitionPermissions(prop.Permissions)
					state.AssignableScopes = pointer.From(prop.AssignableScopes)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenRoleDefinitionPermissions(input *[]roledefinitions.Permission) []PermissionModel {
	permissions := make([]PermissionModel, 0)
	if input == nil {
		return permissions
	}

	for _, permission := range *input {
		permissions = append(permissions, PermissionModel{
			Actions:        pointer.From(permission.Actions),
			DataActions:    pointer.From(permission.DataActions),
			NotActions:     pointer.From(permission.NotActions),
			NotDataActions: pointer.From(permission.NotDataActions),
		})
	}

	return permissions
}

func (Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.RoleDefinitionV0ToV1{},
		},
	}
}

func roleDefinitionEventualConsistencyUpdate(ctx context.Context, client *roledefinitions.RoleDefinitionsClient, id roledefinitions.ScopedRoleDefinitionId, updateRequestDate string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return resp, "Failed", err
		}
		if resp.Model == nil {
			return resp, "Failed", fmt.Errorf("`model` was nil")
		}
		if resp.Model.Properties == nil {
			return resp, "Failed", fmt.Errorf("`properties` was nil")
		}
		if resp.Model.Properties.CreatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.CreatedOn` was nil")
		}

		if resp.Model.Properties.UpdatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.UpdatedOn` was nil")
		}

		updateRequestTime, err := time.Parse(time.RFC3339, updateRequestDate)
		if err != nil {
			return nil, "", fmt.Errorf("parsing time from update request: %+v", err)
		}

		respCreatedOn, err := time.Parse(time.RFC3339, *resp.Model.Properties.CreatedOn)
		if err != nil {
			return nil, "", fmt.Errorf("parsing time for createdOn from update request: %+v", err)
		}

		respUpdatedOn, err := time.Parse(time.RFC3339, *resp.Model.Properties.UpdatedOn)
		if err != nil {
			return nil, "", fmt.Errorf("parsing time for updatedOn from update request: %+v", err)
		}

		if respCreatedOn.Equal(updateRequestTime) {
			// a new role definition is created and eventually (~5s) reconciled
			return resp, "Pending", nil
		}

		if updateRequestTime.After(respUpdatedOn) {
			// The real updated on will be equal or after the time we requested it due to the swap out.
			return resp, "Pending", nil
		}

		return resp, "Updated", nil
	}
}
