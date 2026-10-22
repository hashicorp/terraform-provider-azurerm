// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRulesClient
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

			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("kind") {
				existing.Kind = expandDataCollectionRuleKind(state.Kind)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("data_flow") {
				existing.Properties.DataFlows = expandDataCollectionRuleDataFlows(state.DataFlows)
			}

			if metadata.ResourceData.HasChange("data_sources") {
				dataSource, err := expandDataCollectionRuleDataSources(state.DataSources)
				if err != nil {
					return err
				}
				existing.Properties.DataSources = dataSource
			}

			if metadata.ResourceData.HasChange("data_collection_endpoint_id") {
				if state.DataCollectionEndpointId != "" {
					existing.Properties.DataCollectionEndpointId = pointer.To(state.DataCollectionEndpointId)
				} else {
					existing.Properties.DataCollectionEndpointId = nil
				}
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("destinations") {
				existing.Properties.Destinations = expandDataCollectionRuleDestinations(state.Destinations)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("stream_declaration") {
				existing.Properties.StreamDeclarations = expandDataCollectionRuleStreamDeclarations(state.StreamDeclaration)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
