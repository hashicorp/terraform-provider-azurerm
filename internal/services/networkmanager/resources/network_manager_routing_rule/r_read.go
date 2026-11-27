// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.RoutingRules

			id, err := routingrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema := ManagerRoutingRuleResourceModel{
				Name:             id.RuleName,
				RuleCollectionId: routingrules.NewRuleCollectionID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.RoutingConfigurationName, id.RuleCollectionName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.Description = pointer.From(model.Properties.Description)
					schema.Destination = flattenNetworkManagerRoutingRuleDestination(props.Destination)
					schema.NextHop = flattenNetworkManagerRoutingRuleNextHop(props.NextHop)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&schema)
		},
	}
}

func flattenNetworkManagerRoutingRuleDestination(input routingrules.RoutingRuleRouteDestination) []ManagerRoutingRuleDestination {
	return []ManagerRoutingRuleDestination{
		{
			Address: input.DestinationAddress,
			Type:    string(input.Type),
		},
	}
}

func flattenNetworkManagerRoutingRuleNextHop(input routingrules.RoutingRuleNextHop) []ManagerRoutingRuleNextHop {
	return []ManagerRoutingRuleNextHop{
		{
			Address: pointer.From(input.NextHopAddress),
			Type:    string(input.NextHopType),
		},
	}
}

type ManagerRoutingRuleDestination struct {
	Address string `tfschema:"address"`
	Type    string `tfschema:"type"`
}

type ManagerRoutingRuleNextHop struct {
	Address string `tfschema:"address"`
	Type    string `tfschema:"type"`
}
