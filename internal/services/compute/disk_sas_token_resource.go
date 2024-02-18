// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
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
			_, err := commonids.ParseManagedDiskID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateManagedDiskID,
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

	diskId, err := commonids.ParseManagedDiskID(d.Get("managed_disk_id").(string))
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

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *diskId)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *diskId)
	}
	props := *resp.Model.Properties

	// checking whether disk export SAS URL is active already before creating. If yes, we raise an error
	if *props.DiskState == disks.DiskStateActiveSAS {
		return fmt.Errorf("active SAS Token for Disk Export already exists, cannot create another one %s: %+v", *diskId, err)
	}

	future, err := client.GrantAccess(ctx, *diskId, grantAccessData)
	if err != nil {
		return fmt.Errorf("granting access to %s: %+v", *diskId, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for access to be granted to %s: %+v", *diskId, err)
	}

	lastResponse := future.Poller.LatestResponse()
	if lastResponse == nil {
		return fmt.Errorf("waiting for access to be granted to %s: last response was nil", *diskId)
	}

	var result Result
	if err := lastResponse.Unmarshal(&result); err != nil {
		return fmt.Errorf("retrieving SAS Token for Disk Access %s: %+v", *diskId, err)
	}

	d.SetId(diskId.ID())
	sasToken := result.Properties.Output.AccessSAS
	d.Set("sas_url", sasToken)

	return resourceManagedDiskSasTokenRead(d, meta)
}

func resourceManagedDiskSasTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diskId, err := commonids.ParseManagedDiskID(d.Id())
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

	id, err := commonids.ParseManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	err = client.RevokeAccessThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("revoking access to %s: %+v", *id, err)
	}

	return nil
}
