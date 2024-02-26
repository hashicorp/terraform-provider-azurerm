// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

const roleDefinitionScope = "/"

type Permission struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

type KeyVaultMHSMRoleDefinitionModel struct {
	Name              string       `tfschema:"name"`
	RoleName          string       `tfschema:"role_name"`
	VaultBaseUrl      string       `tfschema:"vault_base_url"`
	Description       string       `tfschema:"description"`
	Permission        []Permission `tfschema:"permission"`
	RoleType          string       `tfschema:"role_type"`
	ResourceManagerId string       `tfschema:"resource_manager_id"`
}

type KeyVaultMHSMRoleDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = KeyVaultMHSMRoleDefinitionResource{}

// Arguments ...
// skip `assignable_scopes` field support as https://github.com/Azure/azure-rest-api-specs/issues/23045
func (k KeyVaultMHSMRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"role_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_base_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"permission": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"not_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_manager_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) ModelObject() interface{} {
	return &KeyVaultMHSMRoleDefinitionModel{}
}

func (k KeyVaultMHSMRoleDefinitionResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_definition"
}

func (k KeyVaultMHSMRoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient

			var model KeyVaultMHSMRoleDefinitionModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			// need a lock for hsm subresource create/update/delete, or API may respond error as below
			// Status=409 Code="Conflict" Message="There was a conflict while trying to delete the role assignment.
			locks.ByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")

			id, err := parse.NewRoleNestedItemID(model.VaultBaseUrl, roleDefinitionScope, parse.RoleDefinitionType, model.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retrieving role definition by name %s: %v", model.Name, err)
				}
				return meta.ResourceRequiresImport(k.ResourceType(), id)
			}

			var param keyvault.RoleDefinitionCreateParameters
			param.Properties = &keyvault.RoleDefinitionProperties{}
			prop := param.Properties
			prop.RoleName = pointer.To(model.RoleName)
			prop.Description = pointer.To(model.Description)
			prop.RoleType = keyvault.RoleTypeCustomRole
			prop.Permissions = expandKeyVaultMHSMRolePermissions(model.Permission)
			prop.AssignableScopes = pointer.To([]keyvault.RoleScope{"/"})

			if _, err = client.CreateOrUpdate(ctx, model.VaultBaseUrl, roleDefinitionScope, model.Name, param); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			// import has no model data but only id
			id, err := parse.RoleNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KeyVaultMHSMRoleDefinitionModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			client := meta.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			result, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
			if err != nil {
				if response.WasNotFound(result.Response.Response) {
					return meta.MarkAsGone(id)
				}
				return err
			}
			model.Name = pointer.From(result.Name)
			model.VaultBaseUrl = id.VaultBaseUrl

			roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(pointer.From(result.ID))
			if err != nil {
				return fmt.Errorf("paring role definition id %s: %v", pointer.From(result.ID), err)
			}
			model.ResourceManagerId = roleID.ID()

			if prop := result.RoleDefinitionProperties; prop != nil {
				model.Description = pointer.ToString(prop.Description)
				model.RoleType = string(prop.RoleType)
				model.RoleName = pointer.From(prop.RoleName)
				model.Permission = flattenKeyVaultMHSMRolePermission(prop.Permissions)
			}

			meta.SetID(id)
			return meta.Encode(&model)
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 10,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient

			var model KeyVaultMHSMRoleDefinitionModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			id, err := parse.NewRoleNestedItemID(model.VaultBaseUrl, roleDefinitionScope, parse.RoleDefinitionType, model.Name)
			if err != nil {
				return err
			}

			locks.ByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")

			existing, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
			if err != nil {
				if response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("not found resource to update: %s", id)
				}
				return fmt.Errorf("retrieving role definition by name %s: %v", model.Name, err)
			}

			var param keyvault.RoleDefinitionCreateParameters
			param.Properties = &keyvault.RoleDefinitionProperties{}
			prop := param.Properties
			prop.RoleName = pointer.To(model.RoleName)
			prop.Description = pointer.To(model.Description)
			prop.RoleType = keyvault.RoleTypeCustomRole
			prop.Permissions = expandKeyVaultMHSMRolePermissions(model.Permission)

			_, err = client.CreateOrUpdate(ctx, model.VaultBaseUrl, roleDefinitionScope, model.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.RoleNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id.ID())

			locks.ByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			if _, err = meta.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient.Delete(ctx, id.VaultBaseUrl, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %+v: %v", id, err)
			}
			return nil
		},
	}
}

func (k KeyVaultMHSMRoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NestedItemId
}

func expandKeyVaultMHSMRolePermissions(perms []Permission) *[]keyvault.Permission {
	var res []keyvault.Permission
	for _, perm := range perms {
		var dataActions, notDataActions []keyvault.DataAction
		for _, data := range perm.DataActions {
			dataActions = append(dataActions, keyvault.DataAction(data))
		}
		for _, notData := range perm.NotDataActions {
			notDataActions = append(notDataActions, keyvault.DataAction(notData))
		}

		res = append(res, keyvault.Permission{
			Actions:        pointer.To(perm.Actions),
			NotActions:     pointer.To(perm.NotActions),
			DataActions:    pointer.To(dataActions),
			NotDataActions: pointer.To(notDataActions),
		})
	}
	return &res
}

func flattenKeyVaultMHSMRolePermission(perms *[]keyvault.Permission) []Permission {
	if perms == nil {
		return make([]Permission, 0)
	}

	var res []Permission
	for _, perm := range *perms {
		var data, notData []string
		for _, item := range pointer.From(perm.DataActions) {
			data = append(data, string(item))
		}
		for _, item := range pointer.From(perm.NotDataActions) {
			notData = append(notData, string(item))
		}

		res = append(res, Permission{
			Actions:        pointer.From(perm.Actions),
			NotActions:     pointer.From(perm.NotActions),
			DataActions:    data,
			NotDataActions: notData,
		})
	}
	return res
}
