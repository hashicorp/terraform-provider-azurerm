package keyvault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultRoleAssignmentModel struct {
	VaultBaseUrl     string `tfschema:"vault_base_url"`
	Name             string `tfschema:"name"`
	Scope            string `tfschema:"scope"`
	RoleDefinitionId string `tfschema:"role_definition_id"`
	PrincipalId      string `tfschema:"principal_id"`
	ResourceId       string `tfschema:"resource_id"`
}

func (k KeyVaultRoleAssignmentModel) ID() string {
	return fmt.Sprintf("%s%s/%s", k.VaultBaseUrl, k.Scope, k.Name)
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
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// https://learn.microsoft.com/en-us/azure/key-vault/managed-hsm/built-in-roles
		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
	return &KeyVaultRoleAssignmentModel{}
}

func (m KeyVaultRoleAssignmentResource) ResourceType() string {
	return "azurerm_key_vault_role_assignment"
}

func (m KeyVaultRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.KeyVault.MHSMRoleAssignClient

			var model KeyVaultRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleAssignmentType, model.Name)
			if err != nil {
				return err
			}
			existing, err := client.Get(ctx, model.VaultBaseUrl, model.Scope, model.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", model.ID(), err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			var param keyvault.RoleAssignmentCreateParameters
			param.Properties = &keyvault.RoleAssignmentProperties{}
			prop := param.Properties
			prop.PrincipalID = pointer.FromString(model.PrincipalId)
			prop.RoleDefinitionID = pointer.FromString(model.RoleDefinitionId)

			_, err = client.Create(ctx, model.VaultBaseUrl, model.Scope, model.Name, param)
			if err != nil {
				return fmt.Errorf("creating %s: %v", model.ID(), err)
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
			client := meta.Client.KeyVault.MHSMRoleAssignClient

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

			var model KeyVaultRoleAssignmentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			prop := result.Properties
			model.Scope = string(prop.Scope)
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

			if _, err := meta.Client.KeyVault.MHSMRoleAssignClient.Delete(ctx, id.VaultBaseUrl, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}
			return nil
		},
	}
}

func (m KeyVaultRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.MHSMNestedItemId
}
