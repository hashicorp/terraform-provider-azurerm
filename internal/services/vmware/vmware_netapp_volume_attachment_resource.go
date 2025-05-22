// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetappFileVolumeAttachment struct {
	Name            string `tfschema:"name"`
	NetAppVolumeId  string `tfschema:"netapp_volume_id"`
	VmwareClusterId string `tfschema:"vmware_cluster_id"`
}

type NetappFileVolumeAttachmentResource struct{}

func (r NetappFileVolumeAttachmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"netapp_volume_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vmware_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterID,
		},
	}
}

func (r NetappFileVolumeAttachmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetappFileVolumeAttachmentResource) ResourceType() string {
	return "azurerm_vmware_netapp_volume_attachment"
}

func (r NetappFileVolumeAttachmentResource) ModelObject() interface{} {
	return &NetappFileVolumeAttachment{}
}

func (r NetappFileVolumeAttachmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastores.ValidateDataStoreID
}

func (r NetappFileVolumeAttachmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Infof("Decoding state...")
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
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := datastores.Datastore{
				Name: utils.String(state.Name),
				Properties: &datastores.DatastoreProperties{
					NetAppVolume: &datastores.NetAppVolume{
						Id: state.NetAppVolumeId,
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r NetappFileVolumeAttachmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Vmware.DataStoreClient
			id, err := datastores.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			clusterId := datastores.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName, id.ClusterName)

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
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

func (r NetappFileVolumeAttachmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Vmware.DataStoreClient
			id, err := datastores.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
