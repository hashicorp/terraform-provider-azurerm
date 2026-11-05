// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroupName, err)
	}

	return nil
}
