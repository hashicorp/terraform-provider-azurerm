// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	// conditionally check for nested resources and error if they exist
	if meta.(*clients.Client).Features.ResourceGroup.PreventDeletionIfContainsResources {
		// Resource groups sometimes hold on to resource information after the resources have been deleted. We'll retry this check to account for that eventual consistency.
		deletePollerContext, deletePollerCancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer deletePollerCancel()

		pollerType := custompollers.NewResourceGroupPreventDeletePoller(client, *id)
		poller := pollers.NewRetryOnErrorPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow, true)
		if err := poller.PollUntilDone(deletePollerContext); err != nil {
			return err
		}
	}

	// This is not the usual pattern for destroys on go-azure-sdk, however, this functionally the same as the resource
	// worked before refactoring, so behaviour has been maintained. This should be investigated in future and brought
	// in-line if possible.
	if resp, err := client.Delete(ctx, *id, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	} else {
		if err := resp.Poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling deleting %s: %+v", *id, err)
		}
	}

	return nil
}
