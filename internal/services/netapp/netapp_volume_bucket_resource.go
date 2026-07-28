// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppVolumeBucketResource struct{}

var (
	_ sdk.Resource             = NetAppVolumeBucketResource{}
	_ sdk.ResourceWithIdentity = NetAppVolumeBucketResource{}
)

func (r NetAppVolumeBucketResource) Identity() resourceids.ResourceId {
	return &buckets.BucketId{}
}

func (r NetAppVolumeBucketResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeBucketModel{}
}

func (r NetAppVolumeBucketResource) ResourceType() string {
	return "azurerm_netapp_volume_bucket"
}

func (r NetAppVolumeBucketResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r NetAppVolumeBucketResource) Arguments() map[string]*pluginsdk.Schema {
	return netAppBucketResourceCommonArguments()
}

func (r NetAppVolumeBucketResource) Attributes() map[string]*pluginsdk.Schema {
	return netAppBucketResourceCommonAttributes()
}

func (r NetAppVolumeBucketResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppVolumeBucketModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(model.VolumeID)
			if err != nil {
				return fmt.Errorf("parsing volume id %s: %+v", model.VolumeID, err)
			}

			id := buckets.NewBucketID(subscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
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
					Permissions: pointer.To(buckets.BucketPermissions(model.Permissions)),
					FileSystemUser: &buckets.FileSystemUser{
						NfsUser:  expandNetAppBucketNfsUser(model.FileSystemNfsUser),
						CifsUser: expandNetAppBucketCifsUser(model.FileSystemCifsUsername),
					},
					AkvDetails: expandNetAppBucketAkvDetails(model.KeyVault),
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

func (r NetAppVolumeBucketResource) Read() sdk.ResourceFunc {
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

func (r NetAppVolumeBucketResource) flatten(metadata sdk.ResourceMetaData, id *buckets.BucketId, bucket *buckets.Bucket) error {
	volumeID := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

	model := netAppModels.NetAppVolumeBucketModel{
		Name:     id.BucketName,
		VolumeID: volumeID.ID(),
	}

	if bucket != nil && bucket.Properties != nil {
		props := bucket.Properties

		model.Path = pointer.From(props.Path)
		model.Permissions = string(pointer.From(props.Permissions))
		model.Status = string(pointer.From(props.Status))

		if props.FileSystemUser != nil {
			model.FileSystemNfsUser = flattenNetAppBucketNfsUser(props.FileSystemUser.NfsUser)
			model.FileSystemCifsUsername = flattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
		}
		model.KeyVault = flattenNetAppBucketAkvDetails(props.AkvDetails)

		if props.Server != nil {
			model.ServerIPAddress = pointer.From(props.Server.IPAddress)
			model.ServerCertificateCommonName = pointer.From(props.Server.CertificateCommonName)
			model.ServerCertificateExpiryDate = pointer.From(props.Server.CertificateExpiryDate)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&model)
}

func (r NetAppVolumeBucketResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state netAppModels.NetAppVolumeBucketModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patchProps := &buckets.BucketPatchProperties{}

			if metadata.ResourceData.HasChange("permissions") {
				patchProps.Permissions = pointer.To(buckets.BucketPatchPermissions(state.Permissions))
			}

			if metadata.ResourceData.HasChange("file_system_nfs_user") || metadata.ResourceData.HasChange("file_system_cifs_username") {
				patchProps.FileSystemUser = &buckets.FileSystemUser{
					NfsUser:  expandNetAppBucketNfsUser(state.FileSystemNfsUser),
					CifsUser: expandNetAppBucketCifsUser(state.FileSystemCifsUsername),
				}
			}

			if metadata.ResourceData.HasChange("key_vault") {
				patchProps.AkvDetails = expandNetAppBucketAkvDetails(state.KeyVault)
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

func (r NetAppVolumeBucketResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
