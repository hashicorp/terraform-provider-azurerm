package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDiskExport() *pluginsdk.Resource {

	return &pluginsdk.Resource{
		Create: resourceDiskExportCreate,
		Read:   resourceDiskExportRead,
		Delete: resourceDiskExportDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagedDiskID,
			},

			"duration_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(30, 4294967295),
			},

			"access_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.AccessLevelRead),
					string(compute.AccessLevelWrite),
				}, false),
			},

			"sas_url": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}

}

func resourceDiskExportCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Disk Export.")
	durationInSeconds := int32(d.Get("duration_in_seconds").(int))
	access := compute.AccessLevel(d.Get("access_level").(string))

	diskId, err := parse.ManagedDiskID(d.Get("managed_disk_id").(string))
	if err != nil {
		return err
	}

	grantAccessData := compute.GrantAccessData{
		Access:            access,
		DurationInSeconds: &durationInSeconds,
	}

	resp, err := client.Get(ctx, diskId.ResourceGroup, diskId.DiskName)
	if err != nil {
		return fmt.Errorf("Error retrieving Disk %s: %+v", *diskId, err)
	}

	// checking whether disk export SAS URL is active already before creating. If yes, we raise an error
	if resp.DiskState == "ActiveSAS" {
		return fmt.Errorf("Active SAS Token for Disk Export already exists, it should be imported %s: %+v", *diskId, err)
	}

	future, err := client.GrantAccess(ctx, diskId.ResourceGroup, diskId.DiskName, grantAccessData)
	if err != nil {
		return fmt.Errorf("granting access to %s: %+v", *diskId, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for access to be granted to %s: %+v", *diskId, err)
	}
	read, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("retrieving SAS Token for Disk Access %s: %+v", *diskId, err)
	}
	if read.AccessSAS == nil {
		return fmt.Errorf("retrieving SAS Token for Disk Access %s: SAS was nil", *diskId)
	}

	d.SetId(diskId.ID())
	sasToken := *read.AccessSAS
	d.Set("sas_url", sasToken)

	return resourceDiskExportRead(d, meta)

}

func resourceDiskExportRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diskId, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, diskId.ResourceGroup, diskId.DiskName)
	if err != nil {
		// checking whether disk export SAS URL is active post creation. If no, we raise an error
		if resp.DiskState != "ActiveSAS" {
			log.Printf("[INFO] Disk SAS token %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Disk Export for SAS Token %s (resource group %s): %s", diskId.DiskName, diskId.ResourceGroup, err)
	}

	d.SetId(diskId.ID())

	return nil
}

func resourceDiskExportDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.RevokeAccess(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		return fmt.Errorf("revoking access to %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for revocation of access to %s: %+v", *id, err)
	}

	return nil
}
