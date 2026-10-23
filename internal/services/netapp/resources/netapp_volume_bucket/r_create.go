// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model models.NetAppVolumeBucketModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(model.VolumeID)
			if err != nil {
				return fmt.Errorf("parsing volume id %s: %+v", model.VolumeID, err)
			}

			id := buckets.NewBucketID(subscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			payload := buckets.Bucket{
				Properties: &buckets.BucketProperties{
					Path:        pointer.To(model.Path),
					Permissions: pointer.ToEnum[buckets.BucketPermissions](model.Permissions),
					FileSystemUser: &buckets.FileSystemUser{
						NfsUser:  helpers.ExpandNetAppBucketNfsUser(model.FileSystemNfsUser),
						CifsUser: helpers.ExpandNetAppBucketCifsUser(model.FileSystemCifsUsername),
					},
					AkvDetails: helpers.ExpandNetAppBucketAkvDetails(model.KeyVault),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
