// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DiskPoolManagedDiskAttachmentResource struct{}

var _ sdk.Resource = DiskPoolManagedDiskAttachmentResource{}
var _ sdk.ResourceWithDeprecationAndNoReplacement = DiskPoolManagedDiskAttachmentResource{}

type DiskPoolManagedDiskAttachmentModel struct {
	DiskPoolId string `tfschema:"disk_pool_id"`
	DiskId     string `tfschema:"managed_disk_id"`
}

func (DiskPoolManagedDiskAttachmentResource) DeprecationMessage() string {
	return "The `azurerm_disk_pool_managed_disk_attachment` resource is deprecated and will be removed in v4.0 of the AzureRM Provider."
}

func (d DiskPoolManagedDiskAttachmentResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: diskpools.ValidateDiskPoolID,
		},
		"managed_disk_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: disks.ValidateDiskID,
		},
	}
}

func (d DiskPoolManagedDiskAttachmentResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (d DiskPoolManagedDiskAttachmentResource) ModelObject() interface{} {
	return &DiskPoolManagedDiskAttachmentModel{}
}

func (d DiskPoolManagedDiskAttachmentResource) ResourceType() string {
	return "azurerm_disk_pool_managed_disk_attachment"
}

func (d DiskPoolManagedDiskAttachmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			attachment := DiskPoolManagedDiskAttachmentModel{}
			err := metadata.Decode(&attachment)
			if err != nil {
				return err
			}
			subscriptionId := metadata.Client.Account.SubscriptionId
			poolId, err := diskpools.ParseDiskPoolID(attachment.DiskPoolId)
			if err != nil {
				return err
			}
			if poolId.SubscriptionId != subscriptionId {
				return fmt.Errorf("Disk Pool subscription id %q is different from provider's subscription", poolId.SubscriptionId)
			}
			diskId, err := disks.ParseDiskID(attachment.DiskId)
			if err != nil {
				return err
			}
			locks.ByID(attachment.DiskPoolId)
			defer locks.UnlockByID(attachment.DiskPoolId)
			id := parse.NewDiskPoolManagedDiskAttachmentId(*poolId, *diskId)

			client := metadata.Client.Disks.DiskPoolsClient
			poolResp, err := client.Get(ctx, *poolId)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", *poolId, err)
			}

			inputDisks := make([]diskpools.Disk, 0)
			if poolResp.Model != nil && poolResp.Model.Properties.Disks != nil {
				inputDisks = *poolResp.Model.Properties.Disks
			}
			for _, disk := range inputDisks {
				existedDiskId, err := disks.ParseDiskID(disk.Id)
				if err != nil {
					return fmt.Errorf("error on parsing existing attached disk id %q %+v", disk.Id, err)
				}
				if *existedDiskId == *diskId {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			inputDisks = append(inputDisks, diskpools.Disk{
				Id: diskId.ID(),
			})

			err = client.UpdateThenPoll(ctx, *poolId, diskpools.DiskPoolUpdate{
				Properties: diskpools.DiskPoolUpdateProperties{
					Disks: &inputDisks,
				},
			})
			if err != nil {
				return fmt.Errorf("creation of %q: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (d DiskPoolManagedDiskAttachmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.DiskPoolManagedDiskAttachmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			poolId := id.DiskPoolId
			client := metadata.Client.Disks.DiskPoolsClient

			poolResp, err := client.Get(ctx, poolId)
			if err != nil {
				if response.WasNotFound(poolResp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving disk pool %q error: %+v", id.DiskPoolId, err)
			}
			if poolResp.Model == nil || poolResp.Model.Properties.Disks == nil {
				return metadata.MarkAsGone(id)
			}

			for _, disk := range *poolResp.Model.Properties.Disks {
				if disk.Id == id.ManagedDiskId.ID() {
					m := DiskPoolManagedDiskAttachmentModel{
						DiskPoolId: id.DiskPoolId.ID(),
						DiskId:     id.ManagedDiskId.ID(),
					}
					return metadata.Encode(&m)
				}
			}

			return metadata.MarkAsGone(id)
		},
	}
}

func (d DiskPoolManagedDiskAttachmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diskToDetach := &DiskPoolManagedDiskAttachmentModel{}
			err := metadata.Decode(diskToDetach)
			if err != nil {
				return err
			}
			poolId, err := diskpools.ParseDiskPoolID(diskToDetach.DiskPoolId)
			if err != nil {
				return err
			}
			locks.ByID(diskToDetach.DiskPoolId)
			defer locks.UnlockByID(diskToDetach.DiskPoolId)

			client := metadata.Client.Disks.DiskPoolsClient
			pool, err := client.Get(ctx, *poolId)
			if err != nil {
				return fmt.Errorf("retrieving disk pool %q error %v", diskToDetach.DiskPoolId, err)
			}
			if pool.Model == nil || pool.Model.Properties.Disks == nil {
				return nil
			}
			attachedDisks := *pool.Model.Properties.Disks
			remainingDisks := make([]diskpools.Disk, 0)
			for _, attachedDisk := range attachedDisks {
				if attachedDisk.Id != diskToDetach.DiskId {
					remainingDisks = append(remainingDisks, attachedDisk)
				}
			}

			err = client.UpdateThenPoll(ctx, *poolId, diskpools.DiskPoolUpdate{
				Properties: diskpools.DiskPoolUpdateProperties{
					Disks: &remainingDisks,
				},
			})
			if err != nil {
				return fmt.Errorf("error on deletion of disk pool managed disk attachment %q: %+v", metadata.ResourceData.Id(), err)
			}
			return nil
		},
	}
}

func (d DiskPoolManagedDiskAttachmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DiskPoolManagedDiskAttachment
}
