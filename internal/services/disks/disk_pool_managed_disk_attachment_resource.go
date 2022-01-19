package disks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/sdk/2021-08-01/diskpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DiskPoolManagedDiskAttachmentResource struct{}

var _ sdk.Resource = DiskPoolManagedDiskAttachmentResource{}

type DiskPoolManagedDiskAttachmentModel struct {
	DiskPoolId string `tfschema:"disk_pool_id"`
	DiskId     string `tfschema:"managed_disk_id"`
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
			ValidateFunc: computeValidate.ManagedDiskID,
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
		Timeout: 30 * time.Minute,
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
			diskId, err := computeParse.ManagedDiskID(attachment.DiskId)
			if err != nil {
				return err
			}
			locks.ByID(attachment.DiskPoolId)
			defer locks.UnlockByID(attachment.DiskPoolId)
			id := diskpools.NewDiskPoolManagedDiskAttachmentId(*poolId, *diskId)

			client := metadata.Client.Disks.DiskPoolsClient
			poolResp, err := client.Get(ctx, *poolId)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", *poolId, err)
			}

			disks := make([]diskpools.Disk, 0)
			if poolResp.Model != nil && poolResp.Model.Properties.Disks != nil {
				disks = *poolResp.Model.Properties.Disks
			}
			for _, disk := range disks {
				existedDiskId, err := computeParse.ManagedDiskID(disk.Id)
				if err != nil {
					return fmt.Errorf("error on parsing existing attached disk id %q %+v", disk.Id, err)
				}
				if *existedDiskId == *diskId {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			disks = append(disks, diskpools.Disk{
				Id: diskId.ID(),
			})

			err = client.UpdateThenPoll(ctx, *poolId, diskpools.DiskPoolUpdate{
				Properties: diskpools.DiskPoolUpdateProperties{
					Disks: &disks,
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
			id, err := diskpools.DiskPoolManagedDiskAttachmentID(metadata.ResourceData.Id())
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
		Timeout: 30 * time.Minute,
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
