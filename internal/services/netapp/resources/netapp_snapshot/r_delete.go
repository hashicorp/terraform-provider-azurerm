// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppSnapshotDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.SnapshotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshots.ParseSnapshotID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// The resource NetApp Snapshot depends on the resource NetApp Volume.
	// Although the delete API returns 404 which means the NetApp Snapshot resource has been deleted.
	// Then it tries to immediately delete NetApp Volume but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to re-check status code again.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/11475
	log.Printf("[DEBUG] Waiting for %s to be deleted", id)
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappSnapshotDeleteStateRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
