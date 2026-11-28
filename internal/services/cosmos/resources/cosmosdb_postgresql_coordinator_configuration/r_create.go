// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_coordinator_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLCoordinatorConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.ConfigurationsClient
			clusterId, err := configurations.ParseServerGroupsv2ID(model.ClusterId)
			if err != nil {
				return err
			}

			id := configurations.NewCoordinatorConfigurationID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ServerGroupsv2Name, model.Name)

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			parameters := configurations.ServerConfiguration{
				Properties: &configurations.ServerConfigurationProperties{
					Value: model.Value,
				},
			}
			if err := client.UpdateOnCoordinatorCallbackThenPoll(ctx, id, parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
