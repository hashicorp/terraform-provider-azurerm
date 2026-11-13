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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.RoutingRules
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerRoutingRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ruleCollectionId, err := routingrules.ParseRuleCollectionID(config.RuleCollectionId)
			if err != nil {
				return err
			}

			id := routingrules.NewRuleID(subscriptionId, ruleCollectionId.ResourceGroupName, ruleCollectionId.NetworkManagerName, ruleCollectionId.RoutingConfigurationName, ruleCollectionId.RuleCollectionName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := routingrules.RoutingRule{
				Name: pointer.To(config.Name),
				Properties: &routingrules.RoutingRulePropertiesFormat{
					Description: pointer.To(config.Description),
					Destination: expandNetworkManagerRoutingRuleDestination(config.Destination),
					NextHop:     expandNetworkManagerRoutingRuleNextHop(config.NextHop),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandNetworkManagerRoutingRuleDestination(input []ManagerRoutingRuleDestination) routingrules.RoutingRuleRouteDestination {
	if len(input) == 0 {
		return routingrules.RoutingRuleRouteDestination{}
	}

	v := input[0]
	return routingrules.RoutingRuleRouteDestination{
		DestinationAddress: v.Address,
		Type:               routingrules.RoutingRuleDestinationType(v.Type),
	}
}

func expandNetworkManagerRoutingRuleNextHop(input []ManagerRoutingRuleNextHop) routingrules.RoutingRuleNextHop {
	if len(input) == 0 {
		return routingrules.RoutingRuleNextHop{}
	}

	v := input[0]

	return routingrules.RoutingRuleNextHop{
		NextHopAddress: pointer.To(v.Address),
		NextHopType:    routingrules.RoutingRuleNextHopType(v.Type),
	}
}
