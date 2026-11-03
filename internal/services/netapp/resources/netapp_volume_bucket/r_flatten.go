// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *buckets.BucketId, bucket *buckets.Bucket) error {
	volumeID := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

	model := models.NetAppVolumeBucketModel{
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
