// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultMHSMRoleDefinitionDataSourceModel struct {
	ManagedHSMID      string       `tfschema:"managed_hsm_id"`
	Name              string       `tfschema:"name"`
	RoleName          string       `tfschema:"role_name"`
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

		"managed_hsm_id": {
			Type:         pluginsdk.TypeString,
			ValidateFunc: managedhsms.ValidateManagedHSMID,
			Required:     true,
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
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMRoleDefinitionDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			var managedHsmId *managedhsms.ManagedHSMId
			var endpoint *parse.ManagedHSMDataPlaneEndpoint
			var err error
			if config.ManagedHSMID != "" {
				managedHsmId, err = managedhsms.ParseManagedHSMID(config.ManagedHSMID)
				if err != nil {
					return err
				}
				baseUri, err := metadata.Client.ManagedHSMs.BaseUriForManagedHSM(ctx, *managedHsmId)
				if err != nil {
					return fmt.Errorf("determining the Data Plane Endpoint for %s: %+v", *managedHsmId, err)
				}
				if baseUri == nil {
					return fmt.Errorf("unable to determine the Data Plane Endpoint for %q", *managedHsmId)
				}
				endpoint, err = parse.ManagedHSMEndpoint(*baseUri, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
				}
			}

			scope := keyvault.RoleScopeGlobal
			id := parse.NewManagedHSMDataPlaneRoleDefinitionID(endpoint.ManagedHSMName, endpoint.DomainSuffix, string(scope), config.Name)

			result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if v := pointer.From(result.ID); v != "" {
				roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(v)
				if err != nil {
					return fmt.Errorf("paring role definition id %q: %v", v, err)
				}
				config.ResourceManagerId = roleID.ID()
			}

			if prop := result.RoleDefinitionProperties; prop != nil {
				config.Description = pointer.ToString(prop.Description)
				config.RoleType = string(prop.RoleType)
				config.RoleName = pointer.From(prop.RoleName)

				if prop.AssignableScopes != nil {
					config.AssignableScopes = make([]string, 0)
					for _, r := range *prop.AssignableScopes {
						config.AssignableScopes = append(config.AssignableScopes, string(r))
					}
				}

				config.Permission = flattenKeyVaultMHSMRolePermission(prop.Permissions)
			}

			metadata.SetID(id)
			return metadata.Encode(&config)
		},
	}
}
