// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

type KeyVaultManagedHSMRoleAssignmentModel struct {
	VaultBaseUrl     string `tfschema:"vault_base_url"`
	Name             string `tfschema:"name"`
	Scope            string `tfschema:"scope"`
	RoleDefinitionId string `tfschema:"role_definition_id"`
	PrincipalId      string `tfschema:"principal_id"`
	ResourceId       string `tfschema:"resource_id"`
}

type KeyVaultManagedHSMRoleAssignmentResource struct{}

var _ sdk.Resource = KeyVaultManagedHSMRoleAssignmentResource{}

func (m KeyVaultManagedHSMRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"vault_base_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(/|/keys|/keys/.+)$`), "scope should be one of `/`, `/keys', `/keys/<key_name>`"),
		},

		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roledefinitions.ValidateScopedRoleDefinitionID,
		},

		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) ModelObject() interface{} {
	return &KeyVaultManagedHSMRoleAssignmentModel{}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_assignment"
}

func (m KeyVaultManagedHSMRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient

			var model KeyVaultManagedHSMRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			locks.ByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")

			id, err := parse.NewRoleNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleAssignmentType, model.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, model.VaultBaseUrl, model.Scope, model.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id.ID(), err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			var param keyvault.RoleAssignmentCreateParameters
			param.Properties = &keyvault.RoleAssignmentProperties{
				PrincipalID: pointer.FromString(model.PrincipalId),
				// the role definition id may has '/' prefix, but the api doesn't accept it
				RoleDefinitionID: pointer.FromString(strings.TrimPrefix(model.RoleDefinitionId, "/")),
			}
			if _, err = client.Create(ctx, model.VaultBaseUrl, model.Scope, model.Name, param); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient

			id, err := parse.RoleNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			var model KeyVaultManagedHSMRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			prop := result.Properties
			model.Name = pointer.From(result.Name)
			model.VaultBaseUrl = id.VaultBaseUrl
			model.Scope = id.Scope
			model.PrincipalId = pointer.ToString(prop.PrincipalID)
			model.ResourceId = pointer.ToString(result.ID)
			if roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(pointer.ToString(prop.RoleDefinitionID)); err != nil {
				return fmt.Errorf("parsing role definition id: %v", err)
			} else {
				model.RoleDefinitionId = roleID.ID()
			}

			return meta.Encode(&model)
		},
	}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.RoleNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)

			locks.ByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			if _, err := meta.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient.Delete(ctx, id.VaultBaseUrl, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}
			return nil
		},
	}
}

func (m KeyVaultManagedHSMRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NestedItemId
}
