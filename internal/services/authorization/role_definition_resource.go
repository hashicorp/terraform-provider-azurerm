// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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

func (r RoleDefinitionResource) Read() sdk.ResourceFunc {
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

func (r RoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			stateId, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			id := roledefinitions.NewScopedRoleDefinitionID(stateId.Scope, stateId.RoleID)

			var config RoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", stateId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", stateId)
			}

			model := *existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", stateId)
			}

			props := *model.Properties

			if metadata.ResourceData.HasChange("name") {
				props.RoleName = &config.Name
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = &config.Description
			}

			if metadata.ResourceData.HasChange("permissions") {
				props.Permissions = pointer.To(expandRoleDefinitionPermissions(config.Permissions))
			}

			if metadata.ResourceData.HasChange("assignable_scopes") {
				props.AssignableScopes = pointer.To(expandRoleDefinitionAssignableScopes(config))
			}

			model.Properties = &props

			resp, err := client.CreateOrUpdate(ctx, id, model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", stateId, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("updating %s: model was nil", stateId)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("updating %s: properties was nil", stateId)
			}

			updatedOn := resp.Model.Properties.UpdatedOn
			if updatedOn == nil {
				return fmt.Errorf("updating Role Definition %q (Scope %q): `properties.UpdatedOn` was nil", stateId.RoleID, stateId.Scope)
			}
			if updatedOn == nil {
				return fmt.Errorf("updating %s: `properties.UpdatedOn` was nil", stateId)
			}

			// "Updating" a role definition actually creates a new one and these get consolidated a few seconds later
			// where the "create date" and "update date" match for the newly created record
			// but eventually switch to being the old create date and the new update date
			// ergo we can can for the old create date and the new updated date
			log.Printf("[DEBUG] Waiting for %s to settle down..", stateId)
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
				Refresh:                   roleDefinitionEventualConsistencyUpdate(ctx, client, id, *updatedOn),
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to settle down: %+v", stateId, err)
			}

			return nil
		},
	}
}

func (r RoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			stateId, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			id := roledefinitions.NewScopedRoleDefinitionID(stateId.Scope, stateId.RoleID)

			resp, err := client.Delete(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("deleting %s: %+v", stateId, err)
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
				Refresh:                   roleDefinitionDeleteStateRefreshFunc(ctx, client, id),
				MinTimeout:                10 * time.Second,
				ContinuousTargetOccurence: 20,
				Timeout:                   time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for delete on Role Definition %s to complete", stateId)
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

func roleDefinitionEventualConsistencyUpdate(ctx context.Context, client *roledefinitions.RoleDefinitionsClient, id roledefinitions.ScopedRoleDefinitionId, updateRequestDate string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
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

func roleDefinitionDeleteStateRefreshFunc(ctx context.Context, client *roledefinitions.RoleDefinitionsClient, id roledefinitions.ScopedRoleDefinitionId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}
			return nil, "Error", err
		}
		return "Pending", "Pending", nil
	}
}
