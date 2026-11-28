// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_coordinator_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ConfigurationsClient

			id, err := configurations.ParseCoordinatorConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetCoordinator(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CosmosDbPostgreSQLCoordinatorConfigurationModel{
				Name:      id.CoordinatorConfigurationName,
				ClusterId: configurations.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Value = props.Value
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}
