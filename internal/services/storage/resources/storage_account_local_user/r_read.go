// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/localuseroperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.LocalUserOperationGroup
			id, err := localuseroperationgroup.ParseLocalUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.LocalUsersGet(ctx, *id)
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

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&model)
		},
	}
}

func (r Resource) flattenPermissionScopes(input *[]localuseroperationgroup.PermissionScope) []PermissionScopeModel {
	if input == nil {
		return []PermissionScopeModel{}
	}

	output := make([]PermissionScopeModel, 0, len(*input))
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
