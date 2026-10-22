// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionendpoints.ParseDataCollectionEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("kind") {
				existing.Kind = expandDataCollectionEndpointKind(state.Kind)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				existing.Properties.NetworkAcls = &datacollectionendpoints.NetworkRuleSet{
					PublicNetworkAccess: expandDataCollectionEndpointPublicNetworkAccess(state.PublicNetworkAccessEnabled),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
