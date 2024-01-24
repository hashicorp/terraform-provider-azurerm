// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RoleDefinitionResource struct{}

var (
	_ sdk.ResourceWithUpdate         = RoleDefinitionResource{}
	_ sdk.ResourceWithStateMigration = RoleDefinitionResource{}
)

type RoleDefinitionModel struct {
	RoleDefinitionId         string            `tfschema:"role_definition_id"`
	Name                     string            `tfschema:"name"`
	Scope                    string            `tfschema:"scope"`
	Description              string            `tfschema:"description"`
	Permissions              []PermissionModel `tfschema:"permissions"`
	AssignableScopes         []string          `tfschema:"assignable_scopes"`
	RoleDefinitionResourceId string            `tfschema:"role_definition_resource_id"`
}

type PermissionModel struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

func (r RoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringStartsWithOneOf("/subscriptions/", "/providers/Microsoft.Management/managementGroups/"),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// lintignore:XS003
		"permissions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"not_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateScopeID,
			},
		},
	}
}

func (r RoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RoleDefinitionResource) ResourceType() string {
	return "azurerm_role_definition"
}

func (r RoleDefinitionResource) ModelObject() interface{} {
	return &RoleDefinitionModel{}
}

func (r RoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.RoleDefinitionId(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (r RoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleDefinitionsClient

			var config RoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			roleDefinitionId := config.RoleDefinitionId
			if roleDefinitionId == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}
				roleDefinitionId = uuid
			}

			existing, err := client.Get(ctx, config.Scope, roleDefinitionId)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Role Definition ID for %q (Scope %q)", config.Name, config.Scope)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				importID := parse.RoleDefinitionID{
					RoleID: roleDefinitionId,
					Scope:  config.Scope,
				}
				return metadata.ResourceRequiresImport(r.ResourceType(), importID)
			}

			properties := authorization.RoleDefinition{
				RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
					RoleName:         &config.Name,
					Description:      &config.Description,
					RoleType:         pointer.To("CustomRole"),
					Permissions:      pointer.To(expandRoleDefinitionPermissions(config.Permissions)),
					AssignableScopes: pointer.To(expandRoleDefinitionAssignableScopes(config)),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, config.Scope, roleDefinitionId, properties); err != nil {
				return err
			}

			read, err := client.Get(ctx, config.Scope, roleDefinitionId)
			if err != nil {
				return err
			}
			if read.ID == nil || *read.ID == "" {
				return fmt.Errorf("cannot read Role Definition ID for %q (Scope %q)", config.Name, config.Scope)
			}

			parsedId := parse.RoleDefinitionID{
				RoleID:     roleDefinitionId,
				Scope:      config.Scope,
				ResourceID: *read.ID,
			}
			metadata.SetID(parsedId)
			return nil
		},
	}
}

func (r RoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleDefinitionsClient

			id, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.Scope, id.RoleID)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := RoleDefinitionModel{
				Scope:            id.Scope,
				RoleDefinitionId: id.RoleID,
			}

			// The Azure resource id of Role Definition is not as same as the one we used to create it.
			// So we read from the response.
			state.RoleDefinitionResourceId = pointer.From(resp.ID)
			state.Name = pointer.From(resp.RoleName)
			state.Description = pointer.From(resp.Description)
			state.Permissions = flattenRoleDefinitionPermissions(resp.Permissions)
			state.AssignableScopes = pointer.From(resp.AssignableScopes)

			return metadata.Encode(&state)
		},
	}
}

func (r RoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			sdkClient := metadata.Client.Authorization.RoleDefinitionsClient
			client := azuresdkhacks.NewRoleDefinitionsWorkaroundClient(sdkClient)

			id, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config RoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			exisiting, err := client.Get(ctx, id.Scope, id.RoleID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			permissions := []authorization.Permission{}
			if config.Permissions != nil {
				for _, permission := range *exisiting.Permissions {
					permissions = append(permissions, authorization.Permission{
						Actions:        permission.Actions,
						DataActions:    permission.DataActions,
						NotActions:     permission.NotActions,
						NotDataActions: permission.NotDataActions,
					})
				}
			}

			update := authorization.RoleDefinition{
				RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
					RoleName:         exisiting.RoleName,
					Description:      exisiting.Description,
					RoleType:         exisiting.RoleType,
					Permissions:      &permissions,
					AssignableScopes: exisiting.AssignableScopes,
				},
			}

			if metadata.ResourceData.HasChange("name") {
				update.RoleDefinitionProperties.RoleName = &config.Name
			}

			if metadata.ResourceData.HasChange("description") {
				update.RoleDefinitionProperties.Description = &config.Description
			}

			if metadata.ResourceData.HasChange("permissions") {
				update.RoleDefinitionProperties.Permissions = pointer.To(expandRoleDefinitionPermissions(config.Permissions))
			}

			if metadata.ResourceData.HasChange("assignable_scopes") {
				update.RoleDefinitionProperties.AssignableScopes = pointer.To(expandRoleDefinitionAssignableScopes(config))
			}

			resp, err := client.CreateOrUpdate(ctx, id.Scope, id.RoleID, update)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			updatedOn := resp.RoleDefinitionProperties.UpdatedOn
			if updatedOn == nil {
				return fmt.Errorf("updating Role Definition %q (Scope %q): `properties.UpdatedOn` was nil", id.RoleID, id.Scope)
			}
			if updatedOn == nil {
				return fmt.Errorf("updating %s: `properties.UpdatedOn` was nil", id)
			}

			// "Updating" a role definition actually creates a new one and these get consolidated a few seconds later
			// where the "create date" and "update date" match for the newly created record
			// but eventually switch to being the old create date and the new update date
			// ergo we can can for the old create date and the new updated date
			log.Printf("[DEBUG] Waiting for %s to settle down..", id)
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				ContinuousTargetOccurence: 12,
				Delay:                     60 * time.Second,
				MinTimeout:                10 * time.Second,
				Pending:                   []string{"Pending"},
				Target:                    []string{"Updated"},
				Refresh:                   roleDefinitionEventualConsistencyUpdate(ctx, client, *id, *updatedOn),
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to settle down: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleDefinitionsClient

			id, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Delete(ctx, id.Scope, id.RoleID)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil
				}
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// Deletes are not instant and can take time to propagate
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					"Pending",
				},
				Target: []string{
					"Deleted",
					"NotFound",
				},
				Refresh:                   roleDefinitionDeleteStateRefreshFunc(ctx, client, *id),
				MinTimeout:                10 * time.Second,
				ContinuousTargetOccurence: 20,
				Timeout:                   time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for delete on Role Definition %s to complete", id)
			}

			return nil
		},
	}
}

func (RoleDefinitionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.RoleDefinitionV0ToV1{},
		},
	}
}

func roleDefinitionEventualConsistencyUpdate(ctx context.Context, client azuresdkhacks.RoleDefinitionsWorkaroundClient, id parse.RoleDefinitionID, updateRequestDate string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.Scope, id.RoleID)
		if err != nil {
			return resp, "Failed", err
		}
		if resp.RoleDefinitionProperties == nil {
			return resp, "Failed", fmt.Errorf("`properties` was nil")
		}
		if resp.RoleDefinitionProperties.CreatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.CreatedOn` was nil")
		}

		if resp.RoleDefinitionProperties.UpdatedOn == nil {
			return resp, "Failed", fmt.Errorf("`properties.UpdatedOn` was nil")
		}

		updateRequestTime, err := time.Parse(time.RFC3339, updateRequestDate)
		if err != nil {
			return nil, "", fmt.Errorf("parsing time from update request: %+v", err)
		}

		respCreatedOn, err := time.Parse(time.RFC3339, *resp.RoleDefinitionProperties.CreatedOn)
		if err != nil {
			return nil, "", fmt.Errorf("parsing time for createdOn from update request: %+v", err)
		}

		respUpdatedOn, err := time.Parse(time.RFC3339, *resp.RoleDefinitionProperties.UpdatedOn)
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

func expandRoleDefinitionPermissions(input []PermissionModel) []authorization.Permission {
	output := make([]authorization.Permission, 0)
	if len(input) == 0 {
		return output
	}

	for _, v := range input {
		permission := authorization.Permission{}

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

func flattenRoleDefinitionPermissions(input *[]authorization.Permission) []PermissionModel {
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

func roleDefinitionDeleteStateRefreshFunc(ctx context.Context, client *authorization.RoleDefinitionsClient, id parse.RoleDefinitionID) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.Scope, id.RoleID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}
			return nil, "Error", err
		}
		return "Pending", "Pending", nil
	}
}
