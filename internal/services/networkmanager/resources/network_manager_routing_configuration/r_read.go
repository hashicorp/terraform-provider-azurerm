// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagerroutingconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagerRoutingConfigurations

			id, err := networkmanagerroutingconfigurations.ParseRoutingConfigurationID(metadata.ResourceData.Id())
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

			networkManagerId := networkmanagerroutingconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID()
			schema := ManagerRoutingConfigurationResourceModel{
				Name:             id.RoutingConfigurationName,
				NetworkManagerId: networkManagerId,
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				schema.Description = pointer.From(model.Properties.Description)
				schema.RouteTableUsageMode = pointer.FromEnum(model.Properties.RouteTableUsageMode)
			}

			return metadata.Encode(&schema)
		},
	}
}
