package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceManagedDiskExportRevoke() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagedDiskExportCancel,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},
		},
	}
}

func dataSourceManagedDiskExportCancel(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedDiskId := d.Get("managed_disk_id").(string)

	parsedManagedDiskId, err := parse.ManagedDiskID(managedDiskId)
	if err != nil {
		return fmt.Errorf("parsing Managed Disk ID %q: %+v", parsedManagedDiskId.ID(), err)
	}

	diskName := parsedManagedDiskId.DiskName
	resourceGroupName := parsedManagedDiskId.ResourceGroup

	// Request to Revoke Access for Disk
	diskRevokeFuture, err := client.RevokeAccess(ctx, resourceGroupName, diskName)
	if err != nil {
		return fmt.Errorf("Error while revoking access for disk export %q: %+v", parsedManagedDiskId.ID(), err)
	}

	// Wait until the Revoke Request is complete
	diskRevokeFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Revoke access operation failed %q (Resource Group %q): %+v", diskName, resourceGroupName, err)
	}

	d.SetId(time.Now().UTC().String())

	return nil

}
