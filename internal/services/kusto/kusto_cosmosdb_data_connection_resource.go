package kusto

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/dataconnections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"
)

type CosmosDBDataConnectionModel struct {
	Name               string `tfschema:"name"`
	ResourceGroupName  string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	CosmosDbAccountId  string `tfschema:"cosmosdb_account_id"`
	CosmosDbContainer  string `tfschema:"cosmosdb_container"`
	CosmosDbDatabase   string `tfschema:"cosmosdb_database"`
	ClusterName        string `tfschema:"cluster_name"`
	DatabaseName       string `tfschema:"database_name"`
	ManagedIdentityId  string `tfschema:"managed_identity_id"`
	MappingRuleName    string `tfschema:"mapping_rule_name"`
	RetrievalStartDate string `tfschema:"retrieval_start_date"`
	TableName          string `tfschema:"table_name"`
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
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"cosmosdb_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"cosmosdb_container": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cosmosdb_database": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},
		"database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseName,
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
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"retrieval_start_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataconnections.NewDataConnectionID(subscriptionId, model.ResourceGroupName, model.ClusterName, model.DatabaseName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := dataconnections.CosmosDbDataConnectionProperties{
				CosmosDbAccountResourceId: model.CosmosDbAccountId,
				CosmosDbContainer:         model.CosmosDbContainer,
				CosmosDbDatabase:          model.CosmosDbDatabase,
				TableName:                 model.TableName,
				ManagedIdentityResourceId: model.ManagedIdentityId,
			}

			if model.MappingRuleName != "" {
				properties.MappingRuleName = &model.MappingRuleName
			}

			if model.RetrievalStartDate != "" {
				properties.RetrievalStartDate = &model.RetrievalStartDate
			}

			l := location.Normalize(model.Location)

			dataConnection := dataconnections.CosmosDbDataConnection{
				Location:   &l,
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
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CosmosDBDataConnectionModel{
				Name:              id.DataConnectionName,
				ResourceGroupName: id.ResourceGroupName,
				ClusterName:       id.ClusterName,
				DatabaseName:      id.DatabaseName,
			}

			if model := resp.Model; model != nil {
				cosmosDbModel := (*model).(dataconnections.CosmosDbDataConnection)
				state.Location = location.Normalize(*cosmosDbModel.Location)

				if properties := cosmosDbModel.Properties; properties != nil {
					state.CosmosDbAccountId = properties.CosmosDbAccountResourceId
					state.CosmosDbContainer = properties.CosmosDbContainer
					state.CosmosDbDatabase = properties.CosmosDbDatabase
					state.TableName = properties.TableName
					state.ManagedIdentityId = properties.ManagedIdentityResourceId

					if properties.MappingRuleName != nil {
						state.MappingRuleName = *properties.MappingRuleName
					}

					if properties.RetrievalStartDate != nil {
						state.RetrievalStartDate = *properties.RetrievalStartDate
					}
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
