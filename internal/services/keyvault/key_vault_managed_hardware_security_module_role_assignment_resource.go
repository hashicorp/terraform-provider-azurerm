package keyvault

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"

	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type ManagedHSMRoleAssignmentModel struct {
	VaultBaseUrl     string `tfschema:"vault_base_url"`
	Name             string `tfschema:"name"`
	Scope            string `tfschema:"scope"`
	RoleDefinitionId string `tfschema:"role_definition_id"`
	PrincipalId      string `tfschema:"principal_id"`
	ResourceId       string `tfschema:"resource_id"`
}

type KeyVaultRoleAssignmentResource struct{}

var _ sdk.Resource = (*KeyVaultRoleAssignmentResource)(nil)

func (m KeyVaultRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
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

func (m KeyVaultRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m KeyVaultRoleAssignmentResource) ModelObject() interface{} {
	return &ManagedHSMRoleAssignmentModel{}
}

func (m KeyVaultRoleAssignmentResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_assignment"
}

func (m KeyVaultRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.KeyVault.MHSMRoleAssignmentsClient

			var model ManagedHSMRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			locks.ByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(model.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")

			id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleAssignmentType, model.Name)
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
				PrincipalID:      pointer.FromString(model.PrincipalId),
				RoleDefinitionID: pointer.FromString(model.RoleDefinitionId),
			}
			_, err = client.Create(ctx, model.VaultBaseUrl, model.Scope, model.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m KeyVaultRoleAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.KeyVault.MHSMRoleAssignmentsClient

			id, err := parse.MHSMNestedItemID(meta.ResourceData.Id())
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

			var model ManagedHSMRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			prop := result.Properties
			model.Name = pointer.From(result.Name)
			model.VaultBaseUrl = id.VaultBaseUrl
			model.Scope = id.Scope

			model.RoleDefinitionId = pointer.ToString(prop.RoleDefinitionID)
			model.PrincipalId = pointer.ToString(prop.PrincipalID)
			model.ResourceId = pointer.ToString(result.ID)

			return meta.Encode(&model)
		},
	}
}

func (m KeyVaultRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.MHSMNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)

			locks.ByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(id.VaultBaseUrl, "azurerm_key_vault_managed_hardware_security_module")
			if _, err := meta.Client.KeyVault.MHSMRoleAssignmentsClient.Delete(ctx, id.VaultBaseUrl, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}
			return nil
		},
	}
}

func (m KeyVaultRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.MHSMNestedItemId
}
