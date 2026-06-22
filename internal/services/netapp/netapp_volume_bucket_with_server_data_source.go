// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppVolumeBucketWithServerDataSource struct{}

var _ sdk.DataSource = NetAppVolumeBucketWithServerDataSource{}

func (r NetAppVolumeBucketWithServerDataSource) ResourceType() string {
	return "azurerm_netapp_volume_bucket_with_server"
}

func (r NetAppVolumeBucketWithServerDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeBucketWithServerDataSourceModel{}
}

func (r NetAppVolumeBucketWithServerDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r NetAppVolumeBucketWithServerDataSource) Arguments() map[string]*pluginsdk.Schema {
	return netAppBucketDataSourceArguments()
}

func (r NetAppVolumeBucketWithServerDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := netAppBucketDataSourceCommonAttributes()

	attributes["server"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"fqdn": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"certificate_pem": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"on_certificate_conflict_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}

	return attributes
}

func (r NetAppVolumeBucketWithServerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			var state netAppModels.NetAppVolumeBucketWithServerDataSourceModel
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
				state.Permissions = string(pointer.From(props.Permissions))
				state.Status = string(pointer.From(props.Status))

				if props.FileSystemUser != nil {
					state.FileSystemNfsUser = flattenNetAppBucketNfsUser(props.FileSystemUser.NfsUser)
					state.FileSystemCifsUser = flattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
				}
				state.KeyVault = flattenNetAppBucketAkvDetails(props.AkvDetails)
				state.Server = flattenNetAppBucketServer(props.Server)

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
