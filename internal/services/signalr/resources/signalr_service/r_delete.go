// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSignalRServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// @tombuildsstuff: we can't use DeleteThenPoll here since the API returns a 404 on the Future in time
	future, err := client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		if r := future.Poller.LatestResponse(); r == nil || !response.WasNotFound(r.Response) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}
