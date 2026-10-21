// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_management_group_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerManagementGroupConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.NetworkManagerConnections
			managementGroupId, err := parse.ManagementGroupID(model.ManagementGroupId)
			if err != nil {
				return err
			}

			id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID(managementGroupId.Name, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, id)
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

			if _, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, id, managerConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
