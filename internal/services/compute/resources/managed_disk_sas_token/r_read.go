// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk_sas_token

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskSasTokenRead(d *pluginsdk.ResourceData, meta any) error {
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

type Result struct {
	Properties Properties `tfschema:"properties"`
}

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}
