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
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			var config RoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			roleId := config.RoleDefinitionId
			if roleId == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}
				roleId = uuid
			}

			id := roledefinitions.NewScopedRoleDefinitionID(config.Scope, roleId)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Role Definition ID for %q (Scope %q)", config.Name, config.Scope)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					importID := parse.RoleDefinitionID{
						RoleID: roleId,
						Scope:  config.Scope,
					}
					return metadata.ResourceRequiresImport(r.ResourceType(), importID)
				}
			}

			properties := roledefinitions.RoleDefinition{
				Properties: &roledefinitions.RoleDefinitionProperties{
					RoleName:         &config.Name,
					Description:      &config.Description,
					Type:             pointer.To("CustomRole"),
					Permissions:      pointer.To(expandRoleDefinitionPermissions(config.Permissions)),
					AssignableScopes: pointer.To(expandRoleDefinitionAssignableScopes(config)),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return err
			}

			read, err := client.Get(ctx, id)
			if err != nil {
				return err
			}

			if read.Model == nil || read.Model.Id == nil || *read.Model.Id == "" {
				return fmt.Errorf("cannot read Role Definition ID for %q (Scope %q)", config.Name, config.Scope)
			}

			stateId := parse.RoleDefinitionID{
				RoleID:     roleId,
				Scope:      config.Scope,
				ResourceID: *read.Model.Id,
			}
			metadata.SetID(stateId)
			return nil
		},
	}
}

func expandRoleDefinitionPermissions(input []PermissionModel) []roledefinitions.Permission {
	output := make([]roledefinitions.Permission, 0)
	if len(input) == 0 {
		return output
	}

	for _, v := range input {
		permission := roledefinitions.Permission{}

		permission.Actions = &v.Actions
		permission.DataActions = &v.DataActions
		permission.NotActions = &v.NotActions
		permission.NotDataActions = &v.NotDataActions

		output = append(output, permission)
	}

	return output
}

func expandRoleDefinitionAssignableScopes(config RoleDefinitionModel) []string {
	scopes := make([]string, 0)

	if len(config.AssignableScopes) == 0 {
		scopes = append(scopes, config.Scope)
	} else {
		scopes = append(scopes, config.AssignableScopes...)
	}

	return scopes
}
