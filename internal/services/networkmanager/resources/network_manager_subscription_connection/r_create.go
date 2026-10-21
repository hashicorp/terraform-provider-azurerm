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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerSubscriptionConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.NetworkManagerConnections
			subscriptionId, err := commonids.ParseSubscriptionID(model.SubscriptionId)
			if err != nil {
				return err
			}

			id := networkmanagerconnections.NewNetworkManagerConnectionID(subscriptionId.SubscriptionId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			managerConnection := networkmanagerconnections.NetworkManagerConnection{
				Properties: &networkmanagerconnections.NetworkManagerConnectionProperties{},
			}

			if model.Description != "" {
				managerConnection.Properties.Description = &model.Description
			}

			if model.NetworkManagerId != "" {
				managerConnection.Properties.NetworkManagerId = &model.NetworkManagerId
			}

			if _, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, id, managerConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
