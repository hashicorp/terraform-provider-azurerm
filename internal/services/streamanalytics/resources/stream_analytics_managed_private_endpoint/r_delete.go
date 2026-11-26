// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_managed_private_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/custompollers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.EndpointsClient
			id, err := privateendpoints.ParsePrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			pollerType := custompollers.NewStreamAnalyticsManagedPrivateEndpointDeletePoller(client, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}
