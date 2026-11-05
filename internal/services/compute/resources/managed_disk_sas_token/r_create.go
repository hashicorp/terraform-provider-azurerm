// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk_sas_token

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskSasTokenCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	// Normally we would use `CallbackThenPoll` but we need the last response for the SAS access URL
	if !meta.(*clients.Client).Features.PersistIDOnCreateBeforePollingForCompletion {
		d.SetId(diskId.ID())
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for access to be granted to %s: %+v", *diskId, err)
	}

	d.SetId(diskId.ID())

	lastResponse := future.Poller.LatestResponse()
	if lastResponse == nil {
		return fmt.Errorf("waiting for access to be granted to %s: last response was nil", *diskId)
	}

	var result Result
	if err := lastResponse.Unmarshal(&result); err != nil {
		return fmt.Errorf("retrieving SAS Token for Disk Access %s: %+v", *diskId, err)
	}

	d.Set("sas_url", result.Properties.Output.AccessSAS)

	return resourceManagedDiskSasTokenRead(d, meta)
}
