// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ResourceType() string {
	return "azurerm_netapp_volume_bucket"
}

func (r DataSource) ModelObject() any {
	return &models.NetAppVolumeBucketDataSourceModel{}
}

func (r DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return helpers.NetAppBucketDataSourceArguments()
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return helpers.NetAppBucketDataSourceCommonAttributes()
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			var state models.NetAppVolumeBucketDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(state.NetAppVolumeID)
			if err != nil {
				return fmt.Errorf("parsing volume id %s: %+v", state.NetAppVolumeID, err)
			}

			id := buckets.NewBucketID(volumeID.SubscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.NetAppVolumeID = volumeID.ID()

			if resp.Model != nil && resp.Model.Properties != nil {
				props := resp.Model.Properties
				state.Path = pointer.From(props.Path)
				state.Permissions = pointer.FromEnum(props.Permissions)
				state.Status = pointer.FromEnum(props.Status)

				if props.FileSystemUser != nil {
					state.FileSystemNfsUser = helpers.FlattenNetAppBucketNfsUser(props.FileSystemUser.NfsUser)
					state.FileSystemCifsUsername = helpers.FlattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
				}
				state.KeyVault = helpers.FlattenNetAppBucketAkvDetails(props.AkvDetails)

				if props.Server != nil {
					state.ServerIPAddress = pointer.From(props.Server.IPAddress)
					state.ServerCertificateCommonName = pointer.From(props.Server.CertificateCommonName)
					state.ServerCertificateExpiryDate = pointer.From(props.Server.CertificateExpiryDate)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
