package keyvault

import (
	"context"
	"fmt"
	"strings"
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

type Permission struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

func (p *Permission) toSDKDataAction() (pda, pnda *[]keyvault.DataAction) {
	var da, nda = make([]keyvault.DataAction, 0), make([]keyvault.DataAction, 0)
	for _, d := range p.DataActions {
		da = append(da, keyvault.DataAction(d))
	}
	for _, nd := range p.NotDataActions {
		nda = append(nda, keyvault.DataAction(nd))
	}
	pda, pnda = &da, &nda
	return
}

func (p *Permission) loadSDKDataAction(perm keyvault.Permission) Permission {
	p.Actions = pointer.ToSliceOfStrings(perm.Actions)
	p.NotActions = pointer.ToSliceOfStrings(perm.NotActions)
	if perm.DataActions != nil {
		for _, a := range *perm.DataActions {
			p.DataActions = append(p.DataActions, string(a))
		}
	}
	if perm.NotDataActions != nil {
		for _, a := range *perm.NotDataActions {
			p.NotDataActions = append(p.NotDataActions, string(a))
		}
	}
	return *p
}

type KeyVaultRoleDefinitionModel struct {
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

func (k *KeyVaultRoleDefinitionModel) id() string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(k.VaultBaseUrl, "/"), k.Scope, k.Name)
}

func (k *KeyVaultRoleDefinitionModel) ToSDKPermissions() *[]keyvault.Permission {
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

func (k *KeyVaultRoleDefinitionModel) LoadSDKPermission(perms *[]keyvault.Permission) {
	if perms != nil {
		k.Permission = []Permission{}
		for _, p := range *perms {
			k.Permission = append(k.Permission, (&Permission{}).loadSDKDataAction(p))
		}
	}
}

type KeyVaultRoleDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = (*KeyVaultRoleDefinitionResource)(nil)

func (k KeyVaultRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"scope": {
			Type:     pluginsdk.TypeString,
			Default:  "/",
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"/",
			}, false),
		},

		//lintignore:XS003
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

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.RoleScopeGlobal),
					// string(keyvault.RoleScopeKeys),
				}, false),
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (k KeyVaultRoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (k KeyVaultRoleDefinitionResource) ModelObject() interface{} {
	return &KeyVaultRoleDefinitionModel{}
}

func (k KeyVaultRoleDefinitionResource) ResourceType() string {
	return "azurerm_key_vault_role_definition"
}

func (k KeyVaultRoleDefinitionResource) createOrUpdateFunc(isUpdate bool) sdk.ResourceRunFunc {
	return func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
		client := meta.Client.KeyVault.MHSMRoleClient

		var model KeyVaultRoleDefinitionModel
		if err = meta.Decode(&model); err != nil {
			return err
		}

		id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleDefinitionType, model.Name)
		if err != nil {
			return err
		}

		existing, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
		if !utils.ResponseWasNotFound(existing.Response) {
			if err != nil {
				return fmt.Errorf("retreiving role definition by name %s: %v", model.Name, err)
			}
			if !isUpdate {
				return meta.ResourceRequiresImport(k.ResourceType(), id)
			}
		} else if isUpdate {
			return fmt.Errorf("not found resource to update: %s", id)
		}

		var param keyvault.RoleDefinitionCreateParameters
		param.Properties = &keyvault.RoleDefinitionProperties{}
		prop := param.Properties
		prop.RoleName = utils.String(model.RoleName)
		prop.Description = utils.String(model.Description)
		prop.RoleType = keyvault.RoleTypeCustomRole
		prop.Permissions = model.ToSDKPermissions()

		var scopes []keyvault.RoleScope
		for _, role := range model.AssignableScopes {
			scopes = append(scopes, keyvault.RoleScope(role))
		}
		if len(scopes) > 0 {
			prop.AssignableScopes = &scopes
		}

		_, err = client.CreateOrUpdate(ctx, model.VaultBaseUrl, model.Scope, model.Name, param)
		if err != nil {
			return fmt.Errorf("creating %s: %v", model.id(), err)
		}

		meta.SetID(id)
		return nil
	}
}

func (k KeyVaultRoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    k.createOrUpdateFunc(false),
	}
}

func (k KeyVaultRoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			// import has no model data but only id
			id, err := parse.MHSMNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KeyVaultRoleDefinitionModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			client := meta.Client.KeyVault.MHSMRoleClient
			result, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
			if utils.ResponseWasNotFound(result.Response) {
				return meta.MarkAsGone(id)
			}
			if err != nil {
				return err
			}
			model.ResourceId = pointer.From(result.ID)

			if prop := result.RoleDefinitionProperties; prop != nil {
				model.Description = pointer.ToString(prop.Description)
				model.RoleType = string(prop.RoleType)
				model.RoleName = pointer.From(prop.RoleName)

				if prop.AssignableScopes != nil {
					model.AssignableScopes = nil
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

func (k KeyVaultRoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 10,
		Func:    k.createOrUpdateFunc(true),
	}
}

func (k KeyVaultRoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.MHSMNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id.ID())

			if _, err = meta.Client.KeyVault.MHSMRoleClient.Delete(ctx, id.VaultBaseUrl, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %+v: %v", id, err)
			}
			return nil
		},
	}
}

func (k KeyVaultRoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.MHSMNestedItemId
}
