// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_netapp_volume_attachment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Vmware.DataStoreClient
			id, err := datastores.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			clusterId := datastores.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName, id.ClusterName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var netAppVolumeId string
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.NetAppVolume != nil {
						netAppVolumeId = props.NetAppVolume.Id
					}
				}
			}
			return metadata.Encode(&NetappFileVolumeAttachment{
				Name:            id.DataStoreName,
				NetAppVolumeId:  netAppVolumeId,
				VmwareClusterId: clusterId.ID(),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

type NetappFileVolumeAttachment struct {
	Name            string `tfschema:"name"`
	NetAppVolumeId  string `tfschema:"netapp_volume_id"`
	VmwareClusterId string `tfschema:"vmware_cluster_id"`
}
