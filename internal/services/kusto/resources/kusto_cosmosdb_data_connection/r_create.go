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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDBDataConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %s: %+v", r.ResourceType(), err)
			}

			client := metadata.Client.Kusto.DataConnectionsClient

			cosmosDbContainerId, err := cosmosdb.ParseContainerID(model.CosmosDbContainerId)
			if err != nil {
				return err
			}

			kustoDatabaseId, err := commonids.ParseKustoDatabaseID(model.DatabaseId)
			if err != nil {
				return err
			}

			// SubscriptionId and ResourceGroupName need to align with the CosmosDB container, and those could be different from the Kusto database
			cosmosDbAccountResourceId := cosmosdb.NewDatabaseAccountID(cosmosDbContainerId.SubscriptionId, cosmosDbContainerId.ResourceGroupName, cosmosDbContainerId.DatabaseAccountName)

			id := dataconnections.NewDataConnectionID(kustoDatabaseId.SubscriptionId, kustoDatabaseId.ResourceGroupName, kustoDatabaseId.KustoClusterName, kustoDatabaseId.KustoDatabaseName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := dataconnections.CosmosDbDataConnectionProperties{
				CosmosDbAccountResourceId: cosmosDbAccountResourceId.ID(),
				CosmosDbContainer:         cosmosDbContainerId.ContainerName,
				CosmosDbDatabase:          cosmosDbContainerId.SqlDatabaseName,
				TableName:                 model.TableName,
				ManagedIdentityResourceId: model.ManagedIdentityId,
			}

			if model.MappingRuleName != "" {
				properties.MappingRuleName = &model.MappingRuleName
			}

			if model.RetrievalStartDate != "" {
				properties.RetrievalStartDate = &model.RetrievalStartDate
			}

			dataConnection := dataconnections.CosmosDbDataConnection{
				Location:   pointer.To(location.Normalize(model.Location)),
				Name:       &model.Name,
				Properties: &properties,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, dataConnection, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
