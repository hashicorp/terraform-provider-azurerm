// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cosmosdb_data_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Kusto.DataConnectionsClient

			id, err := dataconnections.ParseDataConnectionID(metadata.ResourceData.Id())
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

			kustoDatabaseId := commonids.NewKustoDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DatabaseName)

			state := CosmosDBDataConnectionModel{
				Name:       id.DataConnectionName,
				DatabaseId: kustoDatabaseId.ID(),
			}

			if model := resp.Model; model != nil {
				cosmosDbModel := model.(dataconnections.CosmosDbDataConnection)
				state.Location = location.Normalize(*cosmosDbModel.Location)

				if properties := cosmosDbModel.Properties; properties != nil {
					cosmosdbAccountId, err := cosmosdb.ParseDatabaseAccountID(properties.CosmosDbAccountResourceId)
					if err != nil {
						return err
					}
					cosmosDbContainerId := cosmosdb.NewContainerID(cosmosdbAccountId.SubscriptionId, cosmosdbAccountId.ResourceGroupName, cosmosdbAccountId.DatabaseAccountName, properties.CosmosDbDatabase, properties.CosmosDbContainer)
					state.CosmosDbContainerId = cosmosDbContainerId.ID()
					state.TableName = properties.TableName
					state.ManagedIdentityId = properties.ManagedIdentityResourceId
					state.MappingRuleName = pointer.From(properties.MappingRuleName)
					state.RetrievalStartDate = pointer.From(properties.RetrievalStartDate)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
