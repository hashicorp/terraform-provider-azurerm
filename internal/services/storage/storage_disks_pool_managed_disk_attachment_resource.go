package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type DisksPoolManagedDiskAttachmentResource struct{}

var _ sdk.Resource = DisksPoolManagedDiskAttachmentResource{}

type DisksPoolManagedDiskAttachmentModel struct {
	DisksPoolId string `tfschema:"disks_pool_id"`
	DiskId      string `tfschema:"managed_disk_id"`
}

func (d DisksPoolManagedDiskAttachmentResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disks_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageDisksPoolID,
		},
		"managed_disk_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: computeValidate.ManagedDiskID,
		},
	}
}

func (d DisksPoolManagedDiskAttachmentResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (d DisksPoolManagedDiskAttachmentResource) ModelObject() interface{} {
	return &DisksPoolManagedDiskAttachmentModel{}
}

func (d DisksPoolManagedDiskAttachmentResource) ResourceType() string {
	return "azurerm_storage_disks_pool_managed_disk_attachment"
}

func (d DisksPoolManagedDiskAttachmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			attachment := DisksPoolManagedDiskAttachmentModel{}
			err := metadata.Decode(&attachment)
			if err != nil {
				return err
			}
			subscriptionId := metadata.Client.Account.SubscriptionId
			poolId, err := parse.StorageDisksPoolID(attachment.DisksPoolId)
			if err != nil {
				return err
			}
			if poolId.SubscriptionId != subscriptionId {
				return fmt.Errorf("Disks Pool subscription id %q is different from provider's subscription", poolId.SubscriptionId)
			}
			diskId, err := computeParse.ManagedDiskID(attachment.DiskId)
			if err != nil {
				return err
			}
			locks.ByID(attachment.DisksPoolId)
			defer locks.UnlockByID(attachment.DisksPoolId)
			id := parse.NewStorageDisksPoolManagedDiskAttachmentId(poolId.ID(), diskId.ID())

			client := metadata.Client.Storage.DisksPoolsClient
			poolResp, err := client.Get(ctx, poolId.ResourceGroup, poolId.DiskPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", *poolId, err)
			}

			if poolResp.Disks == nil {
				poolResp.Disks = &[]storagepool.Disk{}
			}
			for _, disk := range *poolResp.Disks {
				if disk.ID == nil {
					continue
				}
				existedDiskId, err := computeParse.ManagedDiskID(*disk.ID)
				if err != nil {
					return fmt.Errorf("error on parsing existing attached disk id %q %+v", *disk.ID, err)
				}
				if existedDiskId == diskId {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			disks := append(*poolResp.Disks, storagepool.Disk{
				ID: utils.String(diskId.ID()),
			})

			future, err := client.Update(ctx, poolId.ResourceGroup, poolId.DiskPoolName, storagepool.DiskPoolUpdate{
				DiskPoolUpdateProperties: &storagepool.DiskPoolUpdateProperties{
					Disks: &disks,
				},
			})
			if err != nil {
				return fmt.Errorf("creation of %q: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %q: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (d DisksPoolManagedDiskAttachmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageDisksPoolManagedDiskAttachmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			poolId, _ := parse.StorageDisksPoolID(id.DisksPoolId)
			client := metadata.Client.Storage.DisksPoolsClient

			poolResp, err := client.Get(ctx, poolId.ResourceGroup, poolId.DiskPoolName)
			if err != nil {
				if utils.ResponseWasNotFound(poolResp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving disks pool %q error: %+v", id.DisksPoolId, err)
			}
			if poolResp.DiskPoolProperties == nil || poolResp.DiskPoolProperties.Disks == nil {
				return metadata.MarkAsGone(id)
			}

			for _, disk := range *poolResp.Disks {
				if disk.ID != nil && *disk.ID == id.ManagedDiskId {
					m := DisksPoolManagedDiskAttachmentModel{
						DisksPoolId: id.DisksPoolId,
						DiskId:      id.ManagedDiskId,
					}
					return metadata.Encode(&m)
				}
			}

			return metadata.MarkAsGone(id)
		},
	}
}

func (d DisksPoolManagedDiskAttachmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diskToDetach := &DisksPoolManagedDiskAttachmentModel{}
			err := metadata.Decode(diskToDetach)
			if err != nil {
				return err
			}
			poolId, err := parse.StorageDisksPoolID(diskToDetach.DisksPoolId)
			if err != nil {
				return err
			}
			locks.ByID(diskToDetach.DisksPoolId)
			defer locks.UnlockByID(diskToDetach.DisksPoolId)

			client := metadata.Client.Storage.DisksPoolsClient
			pool, err := client.Get(ctx, poolId.ResourceGroup, poolId.DiskPoolName)
			if err != nil {
				return fmt.Errorf("retrieving disks pool %q error %v", diskToDetach.DisksPoolId, err)
			}
			if pool.Disks == nil {
				return nil
			}
			attachedDisks := *pool.Disks
			remainingDisks := make([]storagepool.Disk, 0)
			for _, attachedDisk := range attachedDisks {
				if utils.NormalizeNilableString(attachedDisk.ID) != diskToDetach.DiskId {
					remainingDisks = append(remainingDisks, attachedDisk)
				}
			}

			future, err := client.Update(ctx, poolId.ResourceGroup, poolId.DiskPoolName, storagepool.DiskPoolUpdate{
				DiskPoolUpdateProperties: &storagepool.DiskPoolUpdateProperties{
					Disks: &remainingDisks,
				},
			})
			if err != nil {
				return fmt.Errorf("error on deletion of disks pool managed disk attachment %q: %+v", metadata.ResourceData.Id(), err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of disks pool managed disk attatchment %q: %+v", metadata.ResourceData.Id(), err)
			}
			return nil
		},
	}
}

func (d DisksPoolManagedDiskAttachmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageDisksPoolManagedDiskAttachment
}
