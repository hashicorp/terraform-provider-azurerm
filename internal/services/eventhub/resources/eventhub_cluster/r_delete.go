// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_cluster

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubClusterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := eventhubsclusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	// The EventHub Cluster can't be deleted until four hours after creation so we'll keep retrying until it can be deleted.
	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutDelete), func() *pluginsdk.RetryError {
		future, err := client.ClustersDelete(ctx, *id)
		if err != nil {
			if strings.Contains(err.Error(), "ClusterMoratoriumInEffect") || response.WasBadRequest(future.HttpResponse) || response.WasStatusCode(future.HttpResponse, http.StatusTooManyRequests) {
				return pluginsdk.RetryableError(fmt.Errorf("expected eventhub cluster to be deleted but was in pending creation state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		if err := future.Poller.PollUntilDone(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		return nil
	})
}
