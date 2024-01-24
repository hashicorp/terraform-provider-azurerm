// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
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
			client := metadata.Client.Authorization.RoleDefinitionsClient

			var config RoleDefinitionDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			defId := config.RoleDefinitionId

			// search by name
			var role authorization.RoleDefinition
			if config.Name != "" {
				// Accounting for eventual consistency
				deadline, ok := ctx.Deadline()
				if !ok {
					return fmt.Errorf("internal error: context had no deadline")
				}
				err := pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
					roleDefinitions, err := client.List(ctx, config.Scope, fmt.Sprintf("roleName eq '%s'", config.Name))
					if err != nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: %+v", err))
					}
					if len(roleDefinitions.Values()) != 1 {
						return pluginsdk.RetryableError(fmt.Errorf("loading Role Definition List: could not find role '%s'", config.Name))
					}
					if roleDefinitions.Values()[0].ID == nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: values[0].ID is nil '%s'", config.Name))
					}

					defId = *roleDefinitions.Values()[0].ID
					role, err = client.GetByID(ctx, defId)
					if err != nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("getting Role Definition by ID %s: %+v", defId, err))
					}
					return nil
				})
				if err != nil {
					return err
				}
			} else {
				var err error
				role, err = client.Get(ctx, config.Scope, defId)
				if err != nil {
					return fmt.Errorf("loading Role Definition: %+v", err)
				}
			}

			state := RoleDefinitionDataSourceModel{
				Scope:            config.Scope,
				RoleDefinitionId: defId,
			}

			state.Name = pointer.From(role.RoleName)
			state.Type = pointer.From(role.Type)
			state.Description = pointer.From(role.Description)
			state.Permissions = flattenDataSourceRoleDefinitionPermissions(role.Permissions)
			state.AssignableScopes = pointer.From(role.AssignableScopes)

			metadata.ResourceData.SetId(*role.ID)
			return metadata.Encode(&state)
		},
	}
}

func flattenDataSourceRoleDefinitionPermissions(input *[]authorization.Permission) []PermissionDataSourceModel {
	permissions := make([]PermissionDataSourceModel, 0)
	if input == nil {
		return permissions
	}

	for _, permission := range *input {
		permissions = append(permissions, PermissionDataSourceModel{
			Actions:        pointer.From(permission.Actions),
			DataActions:    pointer.From(permission.DataActions),
			NotActions:     pointer.From(permission.NotActions),
			NotDataActions: pointer.From(permission.NotDataActions),
		})
	}

	return permissions
}
