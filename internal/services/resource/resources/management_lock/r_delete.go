// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_lock

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementLockDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementlocks.ParseScopedLockID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.DeleteByScope(ctx, *id); err != nil {
		// @tombuildsstuff: this is intentionally here in case the parent is gone, since we're under a scope
		// which isn't ideal (as this logic shouldn't be present for most resources) but should for this one
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	poller := custompollers.NewEventualConsistencyPoller(12, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.GetByScope(pollerCtx, *id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish delete replication", id)
	}

	return nil
}
