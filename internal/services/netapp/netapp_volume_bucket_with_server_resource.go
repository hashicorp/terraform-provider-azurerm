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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppVolumeBucketWithServerResource struct{}

var (
	_ sdk.Resource             = NetAppVolumeBucketWithServerResource{}
	_ sdk.ResourceWithIdentity = NetAppVolumeBucketWithServerResource{}
)

func (r NetAppVolumeBucketWithServerResource) Identity() resourceids.ResourceId {
	return &buckets.BucketId{}
}

func (r NetAppVolumeBucketWithServerResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeBucketWithServerModel{}
}

func (r NetAppVolumeBucketWithServerResource) ResourceType() string {
	return "azurerm_netapp_volume_bucket_with_server"
}

func (r NetAppVolumeBucketWithServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r NetAppVolumeBucketWithServerResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := netAppBucketResourceCommonArguments()

	// When the certificate is sourced from Key Vault the inline `server.0.certificate_pem`
	// must not be supplied, so wire up the mutual exclusivity that only applies when the
	// `server` block exists.
	arguments["key_vault"].ConflictsWith = []string{"server.0.certificate_pem"}

	arguments["server"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"fqdn": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"certificate_pem": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					Sensitive:     true,
					ValidateFunc:  validation.StringIsBase64,
					ConflictsWith: []string{"key_vault"},
				},
				"on_certificate_conflict_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(buckets.OnCertificateConflictActionFail),
					ValidateFunc: validation.StringInSlice(buckets.PossibleValuesForOnCertificateConflictAction(), false),
				},
			},
		},
	}

	return arguments
}

func (r NetAppVolumeBucketWithServerResource) Attributes() map[string]*pluginsdk.Schema {
	return netAppBucketResourceCommonAttributes()
}

func (r NetAppVolumeBucketWithServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppVolumeBucketWithServerModel
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
						CifsUser: expandNetAppBucketCifsUser(model.FileSystemCifsUser),
					},
					AkvDetails: expandNetAppBucketAkvDetails(model.KeyVault),
					Server:     expandNetAppBucketServer(model.Server),
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

func (r NetAppVolumeBucketWithServerResource) Read() sdk.ResourceFunc {
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

func (r NetAppVolumeBucketWithServerResource) flatten(metadata sdk.ResourceMetaData, id *buckets.BucketId, bucket *buckets.Bucket) error {
	volumeID := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

	model := netAppModels.NetAppVolumeBucketWithServerModel{
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
			model.FileSystemCifsUser = flattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
		}
		model.KeyVault = flattenNetAppBucketAkvDetails(props.AkvDetails)
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
			model.Server = []netAppModels.NetAppVolumeBucketServer{{}}
		}
		model.Server[0].CertificatePem = v.(string)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&model)
}

func (r NetAppVolumeBucketWithServerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state netAppModels.NetAppVolumeBucketWithServerModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patchProps := &buckets.BucketPatchProperties{}

			if metadata.ResourceData.HasChange("permissions") {
				patchProps.Permissions = pointer.To(buckets.BucketPatchPermissions(state.Permissions))
			}

			if metadata.ResourceData.HasChange("file_system_nfs_user") || metadata.ResourceData.HasChange("file_system_cifs_user") {
				patchProps.FileSystemUser = &buckets.FileSystemUser{
					NfsUser:  expandNetAppBucketNfsUser(state.FileSystemNfsUser),
					CifsUser: expandNetAppBucketCifsUser(state.FileSystemCifsUser),
				}
			}

			if metadata.ResourceData.HasChange("key_vault") {
				patchProps.AkvDetails = expandNetAppBucketAkvDetails(state.KeyVault)
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

func (r NetAppVolumeBucketWithServerResource) Delete() sdk.ResourceFunc {
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
