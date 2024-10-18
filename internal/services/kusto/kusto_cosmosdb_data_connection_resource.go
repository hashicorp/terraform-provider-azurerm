// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/dataconnections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDBDataConnectionModel struct {
	Name                string `tfschema:"name"`
	Location            string `tfschema:"location"`
	CosmosDbContainerId string `tfschema:"cosmosdb_container_id"`
	DatabaseId          string `tfschema:"kusto_database_id"`
	ManagedIdentityId   string `tfschema:"managed_identity_id"`
	MappingRuleName     string `tfschema:"mapping_rule_name"`
	RetrievalStartDate  string `tfschema:"retrieval_start_date"`
	TableName           string `tfschema:"table_name"`
}

var _ sdk.Resource = CosmosDBDataConnectionResource{}

type CosmosDBDataConnectionResource struct{}

func (r CosmosDBDataConnectionResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataConnectionName,
		},
		"location": commonschema.Location(),
		"cosmosdb_container_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateContainerID,
		},
		"kusto_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoDatabaseID,
		},
		"managed_identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"table_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"mapping_rule_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"retrieval_start_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r CosmosDBDataConnectionResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r CosmosDBDataConnectionResource) ModelObject() interface{} {
	return &CosmosDBDataConnectionModel{}
}

func (r CosmosDBDataConnectionResource) ResourceType() string {
	return "azurerm_kusto_cosmosdb_data_connection"
}

func (r CosmosDBDataConnectionResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, dataConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDBDataConnectionResource) Read() sdk.ResourceFunc {
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

func (r CosmosDBDataConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Kusto.DataConnectionsClient

			id, err := dataconnections.ParseDataConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CosmosDBDataConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnections.ValidateDataConnectionID
}
