// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_subscription_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerSubscriptionConnectionModel{
				Name:           id.NetworkManagerConnectionName,
				SubscriptionId: commonids.NewSubscriptionID(id.SubscriptionId).ID(),
			}

			if properties.ConnectionState != nil {
				state.ConnectionState = string(*properties.ConnectionState)
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.NetworkManagerId != nil {
				state.NetworkManagerId = *properties.NetworkManagerId
			}

			return metadata.Encode(&state)
		},
	}
}
