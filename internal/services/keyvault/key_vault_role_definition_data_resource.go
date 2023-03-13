package keyvault

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultRoleDefinitionDataSourceModel struct {
	Name             string       `tfschema:"name"`
	RoleName         string       `tfschema:"role_name"`
	Scope            string       `tfschema:"scope"`
	VaultBaseUrl     string       `tfschema:"vault_base_url"`
	Description      string       `tfschema:"description"`
	AssignableScopes []string     `tfschema:"assignable_scopes"`
	Permission       []Permission `tfschema:"permission"`
	RoleType         string       `tfschema:"role_type"`
	ResourceId       string       `tfschema:"resource_id"`
}

func (k *KeyVaultRoleDefinitionDataSourceModel) ToSDKPermissions() *[]keyvault.Permission {
	var res []keyvault.Permission
	for _, p := range k.Permission {
		ins := keyvault.Permission{}
		if p.Actions != nil {
			ins.Actions = pointer.FromSliceOfStrings(p.Actions)
		}
		if p.NotActions != nil {
			ins.NotActions = pointer.FromSliceOfStrings(p.NotActions)
		}
		ins.DataActions, ins.NotDataActions = p.toSDKDataAction()
		res = append(res, ins)
	}
	return &res
}

func (k *KeyVaultRoleDefinitionDataSourceModel) LoadSDKPermission(perms *[]keyvault.Permission) {
	if perms != nil {
		k.Permission = []Permission{}
		for _, p := range *perms {
			k.Permission = append(k.Permission, (&Permission{}).loadSDKDataAction(p))
		}
	}
}

type KeyvaultRoleDefinitionDataResource struct{}

var _ sdk.DataSource = (*KeyvaultRoleDefinitionDataResource)(nil)

func (k KeyvaultRoleDefinitionDataResource) Arguments() map[string]*pluginsdk.Schema {
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

		"scope": {
			Type:     pluginsdk.TypeString,
			Default:  "/",
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"/",
			}, false),
		},
	}
}

func (k KeyvaultRoleDefinitionDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"role_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_id": {
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

		//lintignore:XS003
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
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},
	}
}

func (k KeyvaultRoleDefinitionDataResource) ModelObject() interface{} {
	return &KeyVaultRoleDefinitionDataSourceModel{}
}

func (k KeyvaultRoleDefinitionDataResource) ResourceType() string {
	return "azurerm_key_vault_role_definition"
}

func (k KeyvaultRoleDefinitionDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var model KeyVaultRoleDefinitionDataSourceModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleDefinitionType, model.Name)
			if err != nil {
				return err
			}

			client := meta.Client.KeyVault.MHSMRoleClient
			result, err := client.Get(ctx, model.VaultBaseUrl, model.Scope, model.Name)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			model.ResourceId = pointer.From(result.ID)

			if prop := result.RoleDefinitionProperties; prop != nil {
				model.Description = pointer.ToString(prop.Description)
				model.RoleType = string(prop.RoleType)
				model.RoleName = pointer.From(prop.RoleName)

				if prop.AssignableScopes != nil {
					for _, r := range *prop.AssignableScopes {
						model.AssignableScopes = append(model.AssignableScopes, string(r))
					}
				}

				model.LoadSDKPermission(prop.Permissions)
			}

			meta.SetID(id)
			return meta.Encode(&model)
		},
	}
}
