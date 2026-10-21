// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_management_group_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagerConnections

			id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerManagementGroupConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties
			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				}
			}

			if metadata.ResourceData.HasChange("network_manager_id") {
				if model.NetworkManagerId != "" {
					properties.NetworkManagerId = &model.NetworkManagerId
				}
			}

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
