// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_netapp_volume_attachment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NetappFileVolumeAttachment
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Vmware.DataStoreClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vmWareClusterId, err := clusters.ParseClusterID(state.VmwareClusterId)
			if err != nil {
				return fmt.Errorf("parsing vmware cluster id %s err: %+v", state.VmwareClusterId, err)
			}

			id := datastores.NewDataStoreID(subscriptionId, vmWareClusterId.ResourceGroupName, vmWareClusterId.PrivateCloudName, vmWareClusterId.ClusterName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := datastores.Datastore{
				Name: pointer.To(state.Name),
				Properties: &datastores.DatastoreProperties{
					NetAppVolume: &datastores.NetAppVolume{
						Id: state.NetAppVolumeId,
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, input, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
