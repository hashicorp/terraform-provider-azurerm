package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

type DataDiskUpdate struct {
	ManagedDiskID      *string
	NewDiskSize        *int32
	NewEncryptionSetID *string
}

// DeleteManagedDisk takes a list of managed disks and attempts to delete those created from "Empty"
// It is intended to be used with Virtual Machine resources for deletion of Data Disks that are created in-line with
// the VM. This may be parallelised at a later date.
func DeleteManagedDisks(ctx context.Context, client *clients.Client, dataDisks *[]compute.DataDisk) error {
	disksClient := client.Compute.DisksClient
	if dataDisks == nil {
		return nil
	}
	for _, v := range *dataDisks {
		if v.ManagedDisk.ID == nil {
			return fmt.Errorf("could not read disk ID for deletion")
		}
		id, err := parse.ManagedDiskID(*v.ManagedDisk.ID)
		if err != nil {
			return fmt.Errorf("could not parse disk ID for deletion: %+v", err)
		}
		log.Printf("[DEBUG] Attempting to delete %s", *id)
		deleteFuture, err := disksClient.Delete(ctx, id.ResourceGroup, id.DiskName)
		if err != nil {
			return fmt.Errorf("failure deleting Data Disk %q (resource group %q): %+v", id.DiskName, id.ResourceGroup, err)
		}

		if err = deleteFuture.WaitForCompletionRef(ctx, disksClient.Client); err != nil {
			return fmt.Errorf("failure waiting for deletion of Data Disk %q (resource group %q): %+v", id.DiskName, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Successfully deleted %s", *id)
	}

	return nil
}

func UpdateManagedDisks(ctx context.Context, client *compute.DisksClient, dataDiskUpdates []DataDiskUpdate) error {
	if len(dataDiskUpdates) == 0 {
		return nil
	}

	for _, v := range dataDiskUpdates {
		id, err := parse.ManagedDiskID(*v.ManagedDiskID)
		if err != nil {
			return err
		}

		diskUpdate := compute.DiskUpdate{
			DiskUpdateProperties: &compute.DiskUpdateProperties{},
		}

		if v.NewDiskSize != nil {
			diskUpdate.DiskSizeGB = v.NewDiskSize
		}

		if v.NewEncryptionSetID != nil {
			diskUpdate.Encryption = &compute.Encryption{
				DiskEncryptionSetID: v.NewEncryptionSetID,
			}
		}

		updateFuture, err := client.Update(ctx, id.ResourceGroup, id.DiskName, diskUpdate)
		if err != nil {
			return fmt.Errorf("failed updating Data Disk %q (resource group %q): %+v", id.DiskName, id.ResourceGroup, err)
		}

		if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("failed waiting for update of Data Disk %q (resource group %q): %+v", id.DiskName, id.ResourceGroup, err)
		}
	}
	return nil
}
