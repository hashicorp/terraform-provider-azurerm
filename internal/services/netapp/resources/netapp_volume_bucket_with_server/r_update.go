// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket_with_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppVolumeBucketWithServerModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patchProps := &buckets.BucketPatchProperties{}

			if metadata.ResourceData.HasChange("permissions") {
				patchProps.Permissions = pointer.ToEnum[buckets.BucketPatchPermissions](state.Permissions)
			}

			if metadata.ResourceData.HasChanges("file_system_nfs_user", "file_system_cifs_username") {
				patchProps.FileSystemUser = &buckets.FileSystemUser{
					NfsUser:  helpers.ExpandNetAppBucketNfsUser(state.FileSystemNfsUser),
					CifsUser: helpers.ExpandNetAppBucketCifsUser(state.FileSystemCifsUsername),
				}
			}

			if metadata.ResourceData.HasChange("key_vault") {
				patchProps.AkvDetails = helpers.ExpandNetAppBucketAkvDetails(state.KeyVault)
			}

			if metadata.ResourceData.HasChange("server") {
				patchProps.Server = expandNetAppBucketServerPatch(state.Server, metadata.ResourceData)
			}

			payload := buckets.BucketPatch{
				Properties: patchProps,
			}

			if err := client.UpdateThenPoll(ctx, pointer.From(id), payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
