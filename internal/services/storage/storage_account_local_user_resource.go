package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	computevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LocalUserResource struct{}

var _ sdk.ResourceWithUpdate = LocalUserResource{}

type PermissionsModel struct {
	Create bool `tfschema:"create"`
	Delete bool `tfschema:"delete"`
	List   bool `tfschema:"list"`
	Read   bool `tfschema:"read"`
	Write  bool `tfschema:"write"`
}
type PermissionScopeModel struct {
	Permissions  []PermissionsModel `tfschema:"permissions"`
	ResourceName string             `tfschema:"resource_name"`
	Service      string             `tfschema:"service"`
}
type SshAuthorizedKeyModel struct {
	Description string `tfschema:"description"`
	Key         string `tfschema:"key"`
}
type LocalUserModel struct {
	HomeDirectory      string                  `tfschema:"home_directory"`
	Name               string                  `tfschema:"name"`
	Password           string                  `tfschema:"password"`
	PermissionScope    []PermissionScopeModel  `tfschema:"permission_scope"`
	Sid                string                  `tfschema:"sid"`
	SshAuthorizedKey   []SshAuthorizedKeyModel `tfschema:"ssh_authorized_key"`
	SshKeyEnabled      bool                    `tfschema:"ssh_key_enabled"`
	SshPasswordEnabled bool                    `tfschema:"ssh_password_enabled"`
	StorageAccountId   string                  `tfschema:"storage_account_id"`
}

func (r LocalUserResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalUserName,
		},
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageAccountID,
		},
		"ssh_key_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			AtLeastOneOf: []string{"ssh_key_enabled", "ssh_password_enabled"},
		},
		"ssh_password_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			AtLeastOneOf: []string{"ssh_key_enabled", "ssh_password_enabled"},
		},
		"home_directory": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"ssh_authorized_key": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"ssh_key_enabled"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     computevalidate.SSHKey,
						DiffSuppressFunc: compute.SSHKeyDiffSuppress,
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"permission_scope": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"permissions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"read": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"delete": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"list": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"create": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
					"service": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(
							[]string{"blob", "file"},
							false,
						),
					},
					"resource_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r LocalUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sid": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
		"password": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r LocalUserResource) ResourceType() string {
	return "azurerm_storage_account_local_user"
}

func (r LocalUserResource) ModelObject() interface{} {
	return &LocalUserModel{}
}

func (r LocalUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.LocalUserID
}

func (r LocalUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.LocalUsersClient

			var plan LocalUserModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Sanity checks on input
			if plan.SshKeyEnabled != (len(plan.SshAuthorizedKey) != 0) {
				if plan.SshKeyEnabled {
					return fmt.Errorf("`ssh_authorized_key` should be specified when `ssh_key_enabled` is enabled")
				} else {
					return fmt.Errorf("`ssh_authorized_key` should not be specified when `ssh_key_enabled` is disabled")
				}
			}

			accountId, err := parse.StorageAccountID(plan.StorageAccountId)
			if err != nil {
				return err
			}

			id := parse.NewLocalUserID(accountId.SubscriptionId, accountId.ResourceGroup, accountId.Name, plan.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := storage.LocalUser{
				LocalUserProperties: &storage.LocalUserProperties{
					PermissionScopes:  r.expandPermissionScopes(plan.PermissionScope),
					SSHAuthorizedKeys: r.expandSSHAuthorizedKeys(plan.SshAuthorizedKey),
					HasSSHKey:         pointer.To(plan.SshKeyEnabled),
					HasSSHPassword:    pointer.To(plan.SshPasswordEnabled),
				},
			}

			if plan.HomeDirectory != "" {
				params.LocalUserProperties.HomeDirectory = utils.String(plan.HomeDirectory)
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			state := plan
			if plan.SshPasswordEnabled {
				resp, err := client.RegeneratePassword(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
				if err != nil {
					return fmt.Errorf("generating password for %s: %v", id.ID(), err)
				}
				if v := resp.SSHPassword; v != nil {
					state.Password = *v
				}
				if err := metadata.Encode(&state); err != nil {
					return err
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LocalUserResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.LocalUsersClient
			id, err := parse.LocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := LocalUserModel{
				Name:             id.Name,
				StorageAccountId: parse.NewStorageAccountID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName).ID(),
				// Password is only accessible during creation
				Password: state.Password,
				// SshAuthorizedKey is only accessible during creation
				SshAuthorizedKey: state.SshAuthorizedKey,
			}

			if props := existing.LocalUserProperties; props != nil {
				model.PermissionScope = r.flattenPermissionScopes(props.PermissionScopes)
				if props.HomeDirectory != nil {
					model.HomeDirectory = *props.HomeDirectory
				}
				if props.HasSSHKey != nil {
					model.SshKeyEnabled = *props.HasSSHKey
				}
				if props.HasSSHPassword != nil {
					model.SshPasswordEnabled = *props.HasSSHPassword
				}
				if props.Sid != nil {
					model.Sid = *props.Sid
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r LocalUserResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.LocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan LocalUserModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Storage.LocalUsersClient

			params, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if props := params.LocalUserProperties; props != nil {
				if metadata.ResourceData.HasChange("home_directory") {
					if plan.HomeDirectory != "" {
						props.HomeDirectory = &plan.HomeDirectory
					} else {
						props.HomeDirectory = nil
					}
				}
				if metadata.ResourceData.HasChange("permission_scope") {
					props.PermissionScopes = r.expandPermissionScopes(plan.PermissionScope)
				}

				if metadata.ResourceData.HasChange("ssh_key_enabled") {
					props.HasSSHKey = &plan.SshKeyEnabled
				}

				if metadata.ResourceData.HasChange("ssh_password_enabled") {
					props.HasSSHPassword = &plan.SshPasswordEnabled
					if _, isEnabled := metadata.ResourceData.GetChange("ssh_password_enabled"); isEnabled.(bool) {
						state := plan
						resp, err := client.RegeneratePassword(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
						if err != nil {
							return fmt.Errorf("generating password for %s: %v", id.ID(), err)
						}
						if v := resp.SSHPassword; v != nil {
							state.Password = *v
						}
						if err := metadata.Encode(&state); err != nil {
							return err
						}
					}
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, id.Name, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LocalUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.LocalUsersClient

			id, err := parse.LocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.StorageAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalUserResource) expandPermissionScopes(input []PermissionScopeModel) *[]storage.PermissionScope {
	if len(input) == 0 {
		return nil
	}

	var output []storage.PermissionScope

	for _, v := range input {
		// The length constraint is guaranteed by schema
		permissions := v.Permissions[0]
		var permissionStr string
		if permissions.Read {
			permissionStr += "r"
		}
		if permissions.Write {
			permissionStr += "w"
		}
		if permissions.Delete {
			permissionStr += "d"
		}
		if permissions.List {
			permissionStr += "l"
		}
		if permissions.Create {
			permissionStr += "c"
		}

		output = append(output, storage.PermissionScope{
			Permissions:  pointer.To(permissionStr),
			Service:      pointer.To(v.Service),
			ResourceName: pointer.To(v.ResourceName),
		})
	}

	return &output
}

func (r LocalUserResource) flattenPermissionScopes(input *[]storage.PermissionScope) []PermissionScopeModel {
	if input == nil {
		return nil
	}

	var output []PermissionScopeModel

	for _, v := range *input {
		permissions := PermissionsModel{}
		if p := v.Permissions; p != nil {
			// The Storage API's have a history of being case-insensitive, so we case-insensitively check the permission here.
			np := strings.ToLower(*p)
			if strings.Index(np, "r") != -1 {
				permissions.Read = true
			}
			if strings.Index(np, "w") != -1 {
				permissions.Write = true
			}
			if strings.Index(np, "d") != -1 {
				permissions.Delete = true
			}
			if strings.Index(np, "l") != -1 {
				permissions.List = true
			}
			if strings.Index(np, "c") != -1 {
				permissions.Create = true
			}
		}

		var service string
		if v.Service != nil {
			service = *v.Service
		}

		var resourceName string
		if v.ResourceName != nil {
			resourceName = *v.ResourceName
		}

		output = append(output, PermissionScopeModel{
			Permissions:  []PermissionsModel{permissions},
			Service:      service,
			ResourceName: resourceName,
		})
	}

	return output
}

func (r LocalUserResource) expandSSHAuthorizedKeys(input []SshAuthorizedKeyModel) *[]storage.SSHPublicKey {
	if len(input) == 0 {
		return nil
	}

	var output []storage.SSHPublicKey

	for _, v := range input {
		output = append(output, storage.SSHPublicKey{
			Description: pointer.To(v.Description),
			Key:         pointer.To(v.Key),
		})
	}

	return &output
}
