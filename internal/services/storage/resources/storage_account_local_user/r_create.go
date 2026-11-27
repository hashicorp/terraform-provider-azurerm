// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/localuseroperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storageClient.From(metadata.Client).ResourceManager.LocalUserOperationGroup

			var plan LocalUserModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Sanity checks on input
			if plan.SshKeyEnabled != (len(plan.SshAuthorizedKey) != 0) {
				if plan.SshKeyEnabled {
					return errors.New("`ssh_authorized_key` should be specified when `ssh_key_enabled` is enabled")
				} else {
					return errors.New("`ssh_authorized_key` should not be specified when `ssh_key_enabled` is disabled")
				}
			}

			accountId, err := commonids.ParseStorageAccountID(plan.StorageAccountId)
			if err != nil {
				return err
			}

			id := localuseroperationgroup.NewLocalUserID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, plan.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.LocalUsersGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			params := localuseroperationgroup.LocalUser{
				Properties: &localuseroperationgroup.LocalUserProperties{
					PermissionScopes:  r.expandPermissionScopes(plan.PermissionScope),
					SshAuthorizedKeys: r.expandSSHAuthorizedKeys(plan.SshAuthorizedKey),
					HasSshKey:         pointer.To(plan.SshKeyEnabled),
					HasSshPassword:    pointer.To(plan.SshPasswordEnabled),
				},
			}

			if plan.HomeDirectory != "" {
				params.Properties.HomeDirectory = pointer.To(plan.HomeDirectory)
			}

			if _, err = client.LocalUsersCreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			state := plan
			if plan.SshPasswordEnabled {
				resp, err := client.LocalUsersRegeneratePassword(ctx, id)
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

			return nil
		},
	}
}

func (r Resource) expandPermissionScopes(input []PermissionScopeModel) *[]localuseroperationgroup.PermissionScope {
	if len(input) == 0 {
		return nil
	}

	output := make([]localuseroperationgroup.PermissionScope, 0, len(input))
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

		output = append(output, localuseroperationgroup.PermissionScope{
			Permissions:  permissionStr,
			Service:      v.Service,
			ResourceName: v.ResourceName,
		})
	}

	return &output
}

func (r Resource) expandSSHAuthorizedKeys(input []SshAuthorizedKeyModel) *[]localuseroperationgroup.SshPublicKey {
	if len(input) == 0 {
		return nil
	}

	output := make([]localuseroperationgroup.SshPublicKey, 0, len(input))
	for _, v := range input {
		output = append(output, localuseroperationgroup.SshPublicKey{
			Description: pointer.To(v.Description),
			Key:         pointer.To(v.Key),
		})
	}

	return &output
}
