// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := datacollectionendpoints.NewDataCollectionEndpointID(subscriptionId, state.ResourceGroupName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := datacollectionendpoints.DataCollectionEndpointResource{
				Kind:     expandDataCollectionEndpointKind(state.Kind),
				Location: location.Normalize(state.Location),
				Name:     pointer.To(state.Name),
				Properties: &datacollectionendpoints.DataCollectionEndpoint{
					Description: pointer.To(state.Description),
					NetworkAcls: &datacollectionendpoints.NetworkRuleSet{
						PublicNetworkAccess: expandDataCollectionEndpointPublicNetworkAccess(state.PublicNetworkAccessEnabled),
					},
				},
				Tags: tags.Expand(state.Tags),
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandDataCollectionEndpointKind(input string) *datacollectionendpoints.KnownDataCollectionEndpointResourceKind {
	if input == "" {
		return nil
	}

	return pointer.ToEnum[datacollectionendpoints.KnownDataCollectionEndpointResourceKind](input)
}

func expandDataCollectionEndpointPublicNetworkAccess(input bool) *datacollectionendpoints.KnownPublicNetworkAccessOptions {
	var result datacollectionendpoints.KnownPublicNetworkAccessOptions
	if input {
		result = datacollectionendpoints.KnownPublicNetworkAccessOptionsEnabled
	} else {
		result = datacollectionendpoints.KnownPublicNetworkAccessOptionsDisabled
	}
	return &result
}
