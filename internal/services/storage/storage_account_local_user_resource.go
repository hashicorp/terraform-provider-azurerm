// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/localusers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	computevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
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
			ValidateFunc: commonids.ValidateStorageAccountID,
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
							// Replace the string literal with enum once https://github.com/Azure/azure-rest-api-specs/pull/21845 is merged
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
	return localusers.ValidateLocalUserID
}

func (r LocalUserResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if diff.HasChange("ssh_password_enabled") {
				if err := diff.SetNewComputed("password"); err != nil {
					return err
				}
			}
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r LocalUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.LocalUsers

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

			accountId, err := commonids.ParseStorageAccountID(plan.StorageAccountId)
			if err != nil {
				return err
			}

			id := localusers.NewLocalUserID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, plan.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := localusers.LocalUser{
				Properties: &localusers.LocalUserProperties{
					PermissionScopes:  r.expandPermissionScopes(plan.PermissionScope),
					SshAuthorizedKeys: r.expandSSHAuthorizedKeys(plan.SshAuthorizedKey),
					HasSshKey:         pointer.To(plan.SshKeyEnabled),
					HasSshPassword:    pointer.To(plan.SshPasswordEnabled),
				},
			}

			if plan.HomeDirectory != "" {
				params.Properties.HomeDirectory = utils.String(plan.HomeDirectory)
			}

			if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			state := plan
			if plan.SshPasswordEnabled {
				resp, err := client.RegeneratePassword(ctx, id)
				if err != nil {
					return fmt.Errorf("generating password for %s: %v", id.ID(), err)
				}
				if resp.Model == nil {
					return fmt.Errorf("unexpected nil of the generate password response model for %s", id.ID())
				}
				if v := resp.Model.SshPassword; v != nil {
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
			client := metadata.Client.Storage.ResourceManager.LocalUsers
			id, err := localusers.ParseLocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := LocalUserModel{
				Name:             id.LocalUserName,
				StorageAccountId: commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID(),
				// Password is only accessible during creation
				Password: state.Password,
				// SshAuthorizedKey is only accessible during creation, whilst this should be returned as it is not a secret.
				// Opened API issue: https://github.com/Azure/azure-rest-api-specs/issues/21866
				SshAuthorizedKey: state.SshAuthorizedKey,
			}

			if existing.Model != nil && existing.Model.Properties != nil {
				props := existing.Model.Properties
				model.PermissionScope = r.flattenPermissionScopes(props.PermissionScopes)
				if props.HomeDirectory != nil {
					model.HomeDirectory = *props.HomeDirectory
				}
				if props.HasSshKey != nil {
					model.SshKeyEnabled = *props.HasSshKey
				}
				if props.HasSshPassword != nil {
					model.SshPasswordEnabled = *props.HasSshPassword
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
			id, err := localusers.ParseLocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan LocalUserModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Storage.ResourceManager.LocalUsers

			params, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := params.Model
			if model == nil {
				return fmt.Errorf("unexpected nil model for %s", id)
			}

			props := model.Properties
			if props == nil {
				return fmt.Errorf("unexpected nil properties for %s", id)
			}

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
				props.HasSshKey = &plan.SshKeyEnabled
			}

			if metadata.ResourceData.HasChange("ssh_authorized_key") {
				props.SshAuthorizedKeys = r.expandSSHAuthorizedKeys(plan.SshAuthorizedKey)
			}

			if metadata.ResourceData.HasChange("ssh_password_enabled") {
				props.HasSshPassword = &plan.SshPasswordEnabled
				_, isEnabled := metadata.ResourceData.GetChange("ssh_password_enabled")
				state := plan
				if isEnabled.(bool) {
					// If this update is to change the `ssh_password_enabled` from false to true. We'll need to regenerate the password.
					// The previously generated password will be useless, that can't be used to connect (sftp returns permission denied).
					// Also, after `ssh_key_enabled` being set to back true, but without calling the RegeneratePassword(), then if you
					// call GET on the local user again, it returns the `ssh_key_enabled` as false, which indicates that we shall always
					// generate a password when enable the `ssh_key_enabled`.
					resp, err := client.RegeneratePassword(ctx, *id)
					if err != nil {
						return fmt.Errorf("generating password for %s: %v", id.ID(), err)
					}
					if resp.Model == nil {
						return fmt.Errorf("unexpected nil of the generate password response model for %s", id.ID())
					}
					if v := resp.Model.SshPassword; v != nil {
						state.Password = *v
					}
				} else {
					state.Password = ""
				}
				if err := metadata.Encode(&state); err != nil {
					return err
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, localusers.LocalUser{Properties: props}); err != nil {
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
			client := metadata.Client.Storage.ResourceManager.LocalUsers

			id, err := localusers.ParseLocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalUserResource) expandPermissionScopes(input []PermissionScopeModel) *[]localusers.PermissionScope {
	if len(input) == 0 {
		return nil
	}

	var output []localusers.PermissionScope

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

		output = append(output, localusers.PermissionScope{
			Permissions:  permissionStr,
			Service:      v.Service,
			ResourceName: v.ResourceName,
		})
	}

	return &output
}

func (r LocalUserResource) flattenPermissionScopes(input *[]localusers.PermissionScope) []PermissionScopeModel {
	if input == nil {
		return nil
	}

	var output []PermissionScopeModel

	for _, v := range *input {
		permissions := PermissionsModel{}
		// The Storage API's have a history of being case-insensitive, so we case-insensitively check the permission here.
		np := strings.ToLower(v.Permissions)
		if strings.Contains(np, "r") {
			permissions.Read = true
		}
		if strings.Contains(np, "w") {
			permissions.Write = true
		}
		if strings.Contains(np, "d") {
			permissions.Delete = true
		}
		if strings.Contains(np, "l") {
			permissions.List = true
		}
		if strings.Contains(np, "c") {
			permissions.Create = true
		}

		output = append(output, PermissionScopeModel{
			Permissions:  []PermissionsModel{permissions},
			Service:      v.Service,
			ResourceName: v.ResourceName,
		})
	}

	return output
}

func (r LocalUserResource) expandSSHAuthorizedKeys(input []SshAuthorizedKeyModel) *[]localusers.SshPublicKey {
	if len(input) == 0 {
		return nil
	}

	var output []localusers.SshPublicKey

	for _, v := range input {
		output = append(output, localusers.SshPublicKey{
			Description: pointer.To(v.Description),
			Key:         pointer.To(v.Key),
		})
	}

	return &output
}
