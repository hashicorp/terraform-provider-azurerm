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
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseCoordinatorConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)
			defer locks.UnlockByName(id.ServerGroupsv2Name, CosmosDbPostgreSQLClusterResourceName)

			var model CosmosDbPostgreSQLCoordinatorConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("value") {
				parameters := configurations.ServerConfiguration{
					Properties: &configurations.ServerConfigurationProperties{
						Value: model.Value,
					},
				}

				if err := client.UpdateOnCoordinatorThenPoll(ctx, *id, parameters); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
