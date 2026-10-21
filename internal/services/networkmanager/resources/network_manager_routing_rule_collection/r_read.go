// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule_collection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.RoutingRuleCollections

			id, err := routingrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
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

			schema := ManagerRoutingRuleCollectionResourceModel{
				Name:                   id.RuleCollectionName,
				RoutingConfigurationId: routingrulecollections.NewRoutingConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.RoutingConfigurationName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.BgpRoutePropagationEnabled = flattenBgpRoutePropagation(props.DisableBgpRoutePropagation)
					schema.Description = pointer.From(props.Description)

					networkGroupIds, err := flattenNetworkManagerRoutingNetworkGroupIds(props.AppliesTo)
					if err != nil {
						return fmt.Errorf("flattening `network_group_ids`: %+v", err)
					}

					schema.NetworkGroupIds = networkGroupIds
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func flattenBgpRoutePropagation(input *routingrulecollections.DisableBgpRoutePropagation) bool {
	return pointer.From(input) == routingrulecollections.DisableBgpRoutePropagationFalse
}

func flattenNetworkManagerRoutingNetworkGroupIds(inputList []routingrulecollections.NetworkManagerRoutingGroupItem) ([]string, error) {
	outputList := make([]string, 0, len(inputList))

	for _, input := range inputList {
		networkGroupId, err := networkgroups.ParseNetworkGroupIDInsensitively(input.NetworkGroupId)
		if err != nil {
			return nil, err
		}

		outputList = append(outputList, networkGroupId.ID())
	}

	return outputList, nil
}
