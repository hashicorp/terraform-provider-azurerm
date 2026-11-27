// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.SqlDedicatedGatewayClient

			id, err := sqldedicatedgateway.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ServiceGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbSqlDedicatedGatewayModel{
				CosmosDbAccountId: cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName).ID(),
			}

			if props := model.Properties; props != nil {
				existing := props.(sqldedicatedgateway.SqlDedicatedGatewayServiceResourceProperties)

				if existing.InstanceCount != nil {
					state.InstanceCount = *existing.InstanceCount
				}

				if existing.InstanceSize != nil {
					state.InstanceSize = *existing.InstanceSize
				}
			}

			return metadata.Encode(&state)
		},
	}
}
