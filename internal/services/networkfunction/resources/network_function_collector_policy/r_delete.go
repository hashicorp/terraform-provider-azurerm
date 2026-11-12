// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_collector_policy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkFunction.CollectorPoliciesClient

			id, err := collectorpolicies.ParseCollectorPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/25152
			log.Printf("[DEBUG] Waiting for %s to be fully deleted..", *id)
			poller := custompollers.NewEventualConsistencyPoller(20, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, *id)
				return resp.HttpResponse, err
			}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}

			return nil
		},
	}
}
