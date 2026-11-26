// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/custompollers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.ClustersClient
			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			pollerType := custompollers.NewStreamAnalyticsClusterDeletePoller(client, *id)
			poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}
