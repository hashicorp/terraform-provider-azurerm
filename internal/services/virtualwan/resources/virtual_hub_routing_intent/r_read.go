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
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := virtualwanClient.From(metadata.Client).VirtualWANs

			id, err := virtualwans.ParseRoutingIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RoutingIntentGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := VirtualHubRoutingIntentModel{
				Name:         id.RoutingIntentName,
				VirtualHubId: virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RoutingPolicies = flattenRoutingPolicy(props.RoutingPolicies)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenRoutingPolicy(input *[]virtualwans.RoutingPolicy) []RoutingPolicy {
	if input == nil {
		return []RoutingPolicy{}
	}

	result := make([]RoutingPolicy, 0, len(*input))
	for _, v := range *input {
		routingPolicy := RoutingPolicy{
			Destinations: v.Destinations,
			Name:         v.Name,
			NextHop:      v.NextHop,
		}

		result = append(result, routingPolicy)
	}

	return result
}

type RoutingPolicy struct {
	Destinations []string `tfschema:"destinations"`
	Name         string   `tfschema:"name"`
	NextHop      string   `tfschema:"next_hop"`
}
