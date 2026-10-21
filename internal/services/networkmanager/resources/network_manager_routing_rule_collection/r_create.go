// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule_collection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.RoutingRuleCollections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerRoutingRuleCollectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			routingConfigurationId, err := routingrulecollections.ParseRoutingConfigurationID(config.RoutingConfigurationId)
			if err != nil {
				return err
			}

			id := routingrulecollections.NewRuleCollectionID(subscriptionId, routingConfigurationId.ResourceGroupName, routingConfigurationId.NetworkManagerName, routingConfigurationId.RoutingConfigurationName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := routingrulecollections.RoutingRuleCollection{
				Name: pointer.To(config.Name),
				Properties: &routingrulecollections.RoutingRuleCollectionPropertiesFormat{
					AppliesTo:                  expandNetworkManagerRoutingNetworkGroupIds(config.NetworkGroupIds),
					Description:                pointer.To(config.Description),
					DisableBgpRoutePropagation: expandBgpRoutePropagation(config.BgpRoutePropagationEnabled),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandBgpRoutePropagation(input bool) *routingrulecollections.DisableBgpRoutePropagation {
	if input {
		return pointer.To(routingrulecollections.DisableBgpRoutePropagationFalse)
	}
	return pointer.To(routingrulecollections.DisableBgpRoutePropagationTrue)
}

func expandNetworkManagerRoutingNetworkGroupIds(inputList []string) []routingrulecollections.NetworkManagerRoutingGroupItem {
	outputList := make([]routingrulecollections.NetworkManagerRoutingGroupItem, 0, len(inputList))
	for _, v := range inputList {
		output := routingrulecollections.NetworkManagerRoutingGroupItem{
			NetworkGroupId: v,
		}

		outputList = append(outputList, output)
	}

	return outputList
}
