// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.RoutingRules

			id, err := routingrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerRoutingRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := resp.Model

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("destination") {
				parameters.Properties.Destination = expandNetworkManagerRoutingRuleDestination(model.Destination)
			}

			if metadata.ResourceData.HasChange("next_hop") {
				parameters.Properties.NextHop = expandNetworkManagerRoutingRuleNextHop(model.NextHop)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
