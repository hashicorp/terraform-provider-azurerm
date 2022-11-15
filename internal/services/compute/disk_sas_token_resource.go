package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

type Result struct {
	Properties Properties `tfschema:"properties"`
}

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

func resourceManagedDiskSasToken() *pluginsdk.Resource {

	return &pluginsdk.Resource{
		Create: resourceManagedDiskSasTokenCreate,
		Read:   resourceManagedDiskSasTokenRead,
		Delete: resourceManagedDiskSasTokenDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := disks.ParseDiskID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: disks.ValidateDiskID,
			},

			// unable to provide upper value of 4294967295 as it's not comptabile with 32-bit (overflow errors)
			"duration_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(30),
			},

			"access_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.AccessLevelRead),
					string(disks.AccessLevelWrite),
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

func resourceManagedDiskSasTokenCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Disk Export.")
	durationInSeconds := int64(d.Get("duration_in_seconds").(int))
	access := disks.AccessLevel(d.Get("access_level").(string))

	diskId, err := disks.ParseDiskID(d.Get("managed_disk_id").(string))
	if err != nil {
		return err
	}

	grantAccessData := disks.GrantAccessData{
		Access:            access,
		DurationInSeconds: durationInSeconds,
	}

	resp, err := client.Get(ctx, *diskId)
	if err != nil {
		return fmt.Errorf("error retrieving Disk %s: %+v", *diskId, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// checking whether disk export SAS URL is active already before creating. If yes, we raise an error
			if string(*props.DiskState) == "ActiveSAS" {
				return fmt.Errorf("active SAS Token for Disk Export already exists, cannot create another one %s: %+v", *diskId, err)
			}

			future, err := client.GrantAccess(ctx, *diskId, grantAccessData)
			if err != nil {
				return fmt.Errorf("granting access to %s: %+v", *diskId, err)
			}

			if err := future.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for access to be granted to %s: %+v", *diskId, err)
			}

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(future.Poller.HttpResponse.Body)
			if err != nil {
				return err
			}

			var result Result
			err = json.Unmarshal([]byte(buf.String()), &result)
			if err != nil {
				return fmt.Errorf("retrieving SAS Token for Disk Access %s: %+v", *diskId, err)
			}
			if result.Properties.Output.AccessSAS == "" {
				return fmt.Errorf("retrieving SAS Token for Disk Access %s: SAS was nil", *diskId)
			}

			d.SetId(diskId.ID())
			sasToken := result.Properties.Output.AccessSAS
			d.Set("sas_url", sasToken)
		}
	}

	return resourceManagedDiskSasTokenRead(d, meta)

}

func resourceManagedDiskSasTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diskId, err := disks.ParseDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *diskId)
	if err != nil {
		// checking whether disk export SAS URL is active post creation. If no, we raise an error
		if string(*resp.Model.Properties.DiskState) != "ActiveSAS" {
			log.Printf("[INFO] Disk SAS token %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Disk Export for SAS Token %s (resource group %s): %s", diskId.DiskName, diskId.ResourceGroupName, err)
	}

	d.SetId(diskId.ID())

	return nil
}

func resourceManagedDiskSasTokenDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disks.ParseDiskID(d.Id())
	if err != nil {
		return err
	}

	err = client.RevokeAccessThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("revoking access to %s: %+v", *id, err)
	}

	return nil
}
