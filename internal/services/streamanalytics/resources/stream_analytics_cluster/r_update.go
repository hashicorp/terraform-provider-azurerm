// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.ClustersClient

			var state ClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChanges("streaming_capacity", "tags") {
				props := clusters.Cluster{
					Sku: &clusters.ClusterSku{
						Capacity: pointer.To(state.StreamingCapacity),
					},
					Tags: tags.Expand(state.Tags),
				}

				var opts clusters.UpdateOperationOptions
				if err := client.UpdateThenPoll(ctx, *id, props, opts); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}
