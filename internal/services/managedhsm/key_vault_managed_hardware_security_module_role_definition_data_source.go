// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultMHSMRoleDefinitionDataSourceModel struct {
	Name              string       `tfschema:"name"`
	RoleName          string       `tfschema:"role_name"`
	VaultBaseUrl      string       `tfschema:"vault_base_url"`
	Description       string       `tfschema:"description"`
	AssignableScopes  []string     `tfschema:"assignable_scopes"`
	Permission        []Permission `tfschema:"permission"`
	RoleType          string       `tfschema:"role_type"`
	ResourceManagerId string       `tfschema:"resource_manager_id"`
}

type KeyvaultMHSMRoleDefinitionDataSource struct{}

var _ sdk.DataSource = KeyvaultMHSMRoleDefinitionDataSource{}

func (k KeyvaultMHSMRoleDefinitionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"vault_base_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func (k KeyvaultMHSMRoleDefinitionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"role_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_manager_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"permission": {
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
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (k KeyvaultMHSMRoleDefinitionDataSource) ModelObject() interface{} {
	return &KeyVaultMHSMRoleDefinitionDataSourceModel{}
}

func (k KeyvaultMHSMRoleDefinitionDataSource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_definition"
}

func (k KeyvaultMHSMRoleDefinitionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var model KeyVaultMHSMRoleDefinitionDataSourceModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			id, err := parse.NewRoleNestedItemID(model.VaultBaseUrl, roleDefinitionScope, parse.RoleDefinitionType, model.Name)
			if err != nil {
				return err
			}

			client := meta.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			result, err := client.Get(ctx, id.VaultBaseUrl, roleDefinitionScope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return fmt.Errorf("%s does not exist", id)
				}
				return err
			}

			roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(pointer.From(result.ID))
			if err != nil {
				return fmt.Errorf("paring role definition id %s: %v", pointer.From(result.ID), err)
			}
			model.ResourceManagerId = roleID.ID()

			if prop := result.RoleDefinitionProperties; prop != nil {
				model.Description = pointer.ToString(prop.Description)
				model.RoleType = string(prop.RoleType)
				model.RoleName = pointer.From(prop.RoleName)

				if prop.AssignableScopes != nil {
					for _, r := range *prop.AssignableScopes {
						model.AssignableScopes = append(model.AssignableScopes, string(r))
					}
				}

				model.Permission = flattenKeyVaultMHSMRolePermission(prop.Permissions)
			}

			meta.SetID(id)
			return meta.Encode(&model)
		},
	}
}
