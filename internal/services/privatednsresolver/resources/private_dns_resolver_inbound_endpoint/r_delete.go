// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_inbound_endpoint

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			id, err := inboundendpoints.ParseInboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, inboundendpoints.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			poller := custompollers.NewEventualConsistencyPoller(3, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, *id)
				return resp.HttpResponse, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:         1 * time.Minute,
				TargetStatusCode: pointer.To(http.StatusNotFound),
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become deleted: %+v", id, err)
			}

			return nil
		},
	}
}
