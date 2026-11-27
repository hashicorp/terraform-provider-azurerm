// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_routing_intent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VirtualHubRoutingIntentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.VirtualWan.VirtualWANs
			virtualHubId, err := virtualwans.ParseVirtualHubID(model.VirtualHubId)
			if err != nil {
				return err
			}

			id := virtualwans.NewRoutingIntentID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.RoutingIntentGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &virtualwans.RoutingIntent{
				Properties: &virtualwans.RoutingIntentProperties{
					RoutingPolicies: expandRoutingPolicy(model.RoutingPolicies),
				},
			}

			if err := client.RoutingIntentCreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandRoutingPolicy(input []RoutingPolicy) *[]virtualwans.RoutingPolicy {
	result := make([]virtualwans.RoutingPolicy, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		routingPolicy := virtualwans.RoutingPolicy{
			Destinations: v.Destinations,
			Name:         v.Name,
			NextHop:      v.NextHop,
		}

		result = append(result, routingPolicy)
	}

	return &result
}
