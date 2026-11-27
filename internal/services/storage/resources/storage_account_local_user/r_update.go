// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/localuseroperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := localuseroperationgroup.ParseLocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan LocalUserModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := storageClient.From(metadata.Client).ResourceManager.LocalUserOperationGroup

			params, err := client.LocalUsersGet(ctx, *id)
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
				state := plan
				if plan.SshPasswordEnabled {
					// If this update is to change the `ssh_password_enabled` from false to true. We'll need to regenerate the password.
					// The previously generated password will be useless, that can't be used to connect (sftp returns permission denied).
					// Also, after `ssh_key_enabled` being set to back true, but without calling the RegeneratePassword(), then if you
					// call GET on the local user again, it returns the `ssh_key_enabled` as false, which indicates that we shall always
					// generate a password when enable the `ssh_key_enabled`.
					resp, err := client.LocalUsersRegeneratePassword(ctx, *id)
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

			if _, err := client.LocalUsersCreateOrUpdate(ctx, *id, localuseroperationgroup.LocalUser{Properties: props}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
