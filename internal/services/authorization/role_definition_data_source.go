// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RoleDefinitionDataSource struct{}

var _ sdk.DataSource = RoleDefinitionDataSource{}

type RoleDefinitionDataSourceModel struct {
	Name             string                      `tfschema:"name"`
	RoleDefinitionId string                      `tfschema:"role_definition_id"`
	Scope            string                      `tfschema:"scope"`
	Description      string                      `tfschema:"description"`
	Type             string                      `tfschema:"type"`
	Permissions      []PermissionDataSourceModel `tfschema:"permissions"`
	AssignableScopes []string                    `tfschema:"assignable_scopes"`
}

type PermissionDataSourceModel struct {
	Actions          []string `tfschema:"actions"`
	NotActions       []string `tfschema:"not_actions"`
	DataActions      []string `tfschema:"data_actions"`
	NotDataActions   []string `tfschema:"not_data_actions"`
	Condition        string   `tfschema:"condition"`
	ConditionVersion string   `tfschema:"condition_version"`
}

func (a RoleDefinitionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ExactlyOneOf: []string{
				"name",
				"role_definition_id",
			},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"role_definition_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ExactlyOneOf: []string{
				"name",
				"role_definition_id",
			},
			ValidateFunc: validation.IsUUID,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateScopeID,
		},
	}
}

func (a RoleDefinitionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"permissions": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"not_actions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
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

					"condition": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"condition_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (a RoleDefinitionDataSource) ModelObject() interface{} {
	return &RoleDefinitionDataSourceModel{}
}

func (a RoleDefinitionDataSource) ResourceType() string {
	return "azurerm_role_definition"
}

func (a RoleDefinitionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			var config RoleDefinitionDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			defId := config.RoleDefinitionId

			// search by name
			var id roledefinitions.ScopedRoleDefinitionId
			var role roledefinitions.RoleDefinition
			if config.Name != "" {
				// Accounting for eventual consistency
				deadline, ok := ctx.Deadline()
				if !ok {
					return fmt.Errorf("internal error: context had no deadline")
				}
				err := pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
					roleDefinitions, err := client.List(ctx, commonids.NewScopeID(config.Scope), roledefinitions.ListOperationOptions{
						Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", config.Name)),
					})
					if err != nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: %+v", err))
					}
					if roleDefinitions.Model == nil {
						return pluginsdk.RetryableError(fmt.Errorf("loading Role Definition List: model was nil"))
					}
					if len(*roleDefinitions.Model) != 1 {
						return pluginsdk.RetryableError(fmt.Errorf("loading Role Definition List: could not find role '%s'", config.Name))
					}
					if (*roleDefinitions.Model)[0].Name == nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: values[0].NameD is nil '%s'", config.Name))
					}

					defId = *(*roleDefinitions.Model)[0].Id
					id = roledefinitions.NewScopedRoleDefinitionID(config.Scope, *(*roleDefinitions.Model)[0].Name)
					return nil
				})
				if err != nil {
					return err
				}
			} else {
				id = roledefinitions.NewScopedRoleDefinitionID(config.Scope, defId)
			}

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `Model` was nil", id)
			}

			role = *resp.Model

			if role.Id == nil {
				return fmt.Errorf("retrieving %s: `Id` was nil", id)
			}

			state := RoleDefinitionDataSourceModel{
				Scope:            id.Scope,
				RoleDefinitionId: defId,
			}
			if props := role.Properties; props != nil {
				state.Name = pointer.From(props.RoleName)
				state.Type = pointer.From(props.Type)
				state.Description = pointer.From(props.Description)
				state.Permissions = flattenDataSourceRoleDefinitionPermissions(props.Permissions)
				state.AssignableScopes = pointer.From(props.AssignableScopes)
			}

			// The sdk managed id start with two "/" when scope is tenant level (empty).
			// So we use the id from response without parsing and reformatting it.
			// Tracked on https://github.com/hashicorp/pandora/issues/3257
			metadata.ResourceData.SetId(*role.Id)
			return metadata.Encode(&state)
		},
	}
}

func flattenDataSourceRoleDefinitionPermissions(input *[]roledefinitions.Permission) []PermissionDataSourceModel {
	permissions := make([]PermissionDataSourceModel, 0)
	if input == nil {
		return permissions
	}
	for _, permission := range *input {
		permissions = append(permissions, PermissionDataSourceModel{
			Actions:          pointer.From(permission.Actions),
			DataActions:      pointer.From(permission.DataActions),
			NotActions:       pointer.From(permission.NotActions),
			NotDataActions:   pointer.From(permission.NotDataActions),
			Condition:        pointer.From(permission.Condition),
			ConditionVersion: pointer.From(permission.ConditionVersion),
		})
	}
	return permissions
}
