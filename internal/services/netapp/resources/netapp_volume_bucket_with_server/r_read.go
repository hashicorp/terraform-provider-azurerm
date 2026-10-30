// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket_with_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *buckets.BucketId, bucket *buckets.Bucket) error {
	volumeID := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

	model := models.NetAppVolumeBucketWithServerModel{
		Name:     id.BucketName,
		VolumeID: volumeID.ID(),
	}

	if bucket != nil && bucket.Properties != nil {
		props := bucket.Properties

		model.Path = pointer.From(props.Path)
		model.Permissions = pointer.FromEnum(props.Permissions)
		model.Status = pointer.FromEnum(props.Status)

		if props.FileSystemUser != nil {
			model.FileSystemNfsUser = helpers.FlattenNetAppBucketNfsUser(props.FileSystemUser.NfsUser)
			model.FileSystemCifsUsername = helpers.FlattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
		}
		model.KeyVault = helpers.FlattenNetAppBucketAkvDetails(props.AkvDetails)
		model.Server = flattenNetAppBucketServer(props.Server)

		if props.Server != nil {
			model.ServerIPAddress = pointer.From(props.Server.IPAddress)
			model.ServerCertificateCommonName = pointer.From(props.Server.CertificateCommonName)
			model.ServerCertificateExpiryDate = pointer.From(props.Server.CertificateExpiryDate)
		}
	}

	// certificate_pem is never returned by the API; preserve from config/state.
	if v, ok := metadata.ResourceData.GetOk("server.0.certificate_pem"); ok {
		if len(model.Server) == 0 {
			model.Server = []models.NetAppVolumeBucketServer{{}}
		}
		model.Server[0].CertificatePem = v.(string)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&model)
}
