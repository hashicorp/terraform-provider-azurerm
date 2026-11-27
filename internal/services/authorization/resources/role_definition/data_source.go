// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type RoleDefinitionDataSourceModel struct {
	Name                     string                      `tfschema:"name"`
	RoleDefinitionId         string                      `tfschema:"role_definition_id"`
	Scope                    string                      `tfschema:"scope"`
	Description              string                      `tfschema:"description"`
	Type                     string                      `tfschema:"type"`
	Permissions              []PermissionDataSourceModel `tfschema:"permissions"`
	AssignableScopes         []string                    `tfschema:"assignable_scopes"`
	RoleDefinitionResourceId string                      `tfschema:"role_definition_resource_id"`
}

type PermissionDataSourceModel struct {
	Actions          []string `tfschema:"actions"`
	NotActions       []string `tfschema:"not_actions"`
	DataActions      []string `tfschema:"data_actions"`
	NotDataActions   []string `tfschema:"not_data_actions"`
	Condition        string   `tfschema:"condition"`
	ConditionVersion string   `tfschema:"condition_version"`
}

func (a DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The Name of either a built-in or custom Role Definition.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{
				"name",
				"role_definition_id",
			},
		},

		"role_definition_id": {
			Description:  "The ID of the Role Definition as a UUID/GUID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IsUUID,
			ExactlyOneOf: []string{
				"name",
				"role_definition_id",
			},
		},

		"scope": {
			Description:  "The Scope at which the Custom Role Definition exists.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateScopeID,
		},
	}
}

func (a DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Description: "The Description of the built-in Role.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"type": {
			Description: "The Type of the Role.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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
						Computed: true,
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
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
			Description: "One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"role_definition_resource_id": {
			Description: "The Azure Resource Manager ID for the resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (a DataSource) ModelObject() any {
	return &RoleDefinitionDataSourceModel{}
}

func (a DataSource) ResourceType() string {
	return "azurerm_role_definition"
}

func (a DataSource) Description() string {
	return "Use this data source to access information about an existing Role Definition."
}

func (a DataSource) Read() sdk.ResourceFunc {
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
					return errors.New("internal error: context had no deadline")
				}
				if err := pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
					roleDefinitions, err := client.List(ctx, commonids.NewScopeID(config.Scope), roledefinitions.ListOperationOptions{
						Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", config.Name)),
					})
					if err != nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: %+v", err))
					}
					if roleDefinitions.Model == nil {
						return pluginsdk.RetryableError(errors.New("loading Role Definition List: model was nil"))
					}
					if len(*roleDefinitions.Model) != 1 {
						return pluginsdk.RetryableError(fmt.Errorf("loading Role Definition List: could not find role '%s'", config.Name))
					}
					if (*roleDefinitions.Model)[0].Name == nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("loading Role Definition List: values[0].NameD is nil '%s'", config.Name))
					}

					defId = *(*roleDefinitions.Model)[0].Name
					id = roledefinitions.NewScopedRoleDefinitionID(config.Scope, defId)

					return nil
				}); err != nil {
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
				Scope:                    id.Scope,
				RoleDefinitionId:         defId,
				RoleDefinitionResourceId: pointer.From(role.Id),
			}

			if props := role.Properties; props != nil {
				state.Name = pointer.From(props.RoleName)
				state.Type = pointer.From(props.Type)
				state.Description = pointer.From(props.Description)
				state.Permissions = flattenDataSourceRoleDefinitionPermissions(props.Permissions)
				state.AssignableScopes = pointer.From(props.AssignableScopes)
			}

			metadata.ResourceData.SetId(*role.Id) // azignore:AZR001 - tenant level (empty) scope IDs start with "//" which the ID parser cannot round-trip, so the raw response ID is used (hashicorp/pandora#3257)
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
