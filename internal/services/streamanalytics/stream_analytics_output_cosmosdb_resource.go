package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2020-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cosmosParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	cosmosValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OutputCosmosDBResource struct{}

var _ sdk.ResourceWithCustomImporter = OutputCosmosDBResource{}

type OutputCosmosDBResourceModel struct {
	Name               string `tfschema:"name"`
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_id"`
	AccountKey         string `tfschema:"cosmosdb_account_key"`
	Database           string `tfschema:"cosmosdb_sql_database_id"`
	ContainerName      string `tfschema:"container_name"`
	DocumentID         string `tfschema:"document_id"`
	PartitionKey       string `tfschema:"partition_key"`
}

func (r OutputCosmosDBResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StreamingJobID,
		},

		"cosmosdb_account_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cosmosdb_sql_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: cosmosValidate.SqlDatabaseID,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"document_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partition_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r OutputCosmosDBResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r OutputCosmosDBResource) ModelObject() interface{} {
	return &OutputCosmosDBResourceModel{}
}

func (r OutputCosmosDBResource) ResourceType() string {
	return "azurerm_stream_analytics_output_cosmosdb"
}

func (r OutputCosmosDBResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputCosmosDBResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.OutputsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			streamingJobId, err := parse.StreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}
			id := parse.NewOutputID(subscriptionId, streamingJobId.ResourceGroup, streamingJobId.Name, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			databaseId, err := cosmosParse.SqlDatabaseID(model.Database)
			if err != nil {
				return err
			}

			documentDbOutputProps := &streamanalytics.DocumentDbOutputDataSourceProperties{
				AccountID:             utils.String(databaseId.DatabaseAccountName),
				AccountKey:            utils.String(model.AccountKey),
				Database:              utils.String(databaseId.Name),
				CollectionNamePattern: utils.String(model.ContainerName),
				DocumentID:            utils.String(model.DocumentID),
				PartitionKey:          utils.String(model.PartitionKey),
			}

			props := streamanalytics.Output{
				Name: utils.String(model.Name),
				OutputProperties: &streamanalytics.OutputProperties{
					Datasource: &streamanalytics.DocumentDbOutputDataSource{
						DocumentDbOutputDataSourceProperties: documentDbOutputProps,
						Type:                                 streamanalytics.TypeBasicOutputDataSourceTypeMicrosoftStorageDocumentDB,
					},
				},
			}

			if _, err := client.CreateOrReplace(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, "", ""); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func (r OutputCosmosDBResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient

			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if props := resp.OutputProperties; props != nil && props.Datasource != nil {
				v, ok := props.Datasource.AsDocumentDbOutputDataSource()
				if !ok {
					return fmt.Errorf("converting output data source to a document DB output: %+v", err)
				}

				streamingJobId := parse.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingjobName)
				state := OutputCosmosDBResourceModel{
					Name:               id.Name,
					StreamAnalyticsJob: streamingJobId.ID(),
				}

				state.AccountKey = metadata.ResourceData.Get("cosmosdb_account_key").(string)

				databaseId := cosmosParse.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, *v.AccountID, *v.Database)
				state.Database = databaseId.ID()

				if v.CollectionNamePattern != nil {
					state.ContainerName = *v.CollectionNamePattern
				}

				if v.DocumentID != nil {
					state.DocumentID = *v.DocumentID
				}

				if v.PartitionKey != nil {
					state.PartitionKey = *v.PartitionKey
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r OutputCosmosDBResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.Name); err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r OutputCosmosDBResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.OutputID
}

func (r OutputCosmosDBResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := parse.OutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputCosmosDBResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			databaseId, err := cosmosParse.SqlDatabaseID(state.Database)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChangesExcept("name", "stream_analytics_job_id") {
				props := streamanalytics.Output{
					OutputProperties: &streamanalytics.OutputProperties{
						Datasource: streamanalytics.DocumentDbOutputDataSource{
							Type: streamanalytics.TypeBasicOutputDataSourceTypeMicrosoftStorageDocumentDB,
							DocumentDbOutputDataSourceProperties: &streamanalytics.DocumentDbOutputDataSourceProperties{
								AccountKey:            &state.AccountKey,
								Database:              &databaseId.Name,
								CollectionNamePattern: &state.ContainerName,
								DocumentID:            &state.DocumentID,
								PartitionKey:          &state.PartitionKey,
							},
						},
					},
				}
				if _, err := client.Update(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, ""); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r OutputCosmosDBResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := parse.OutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
		if err != nil || resp.OutputProperties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.OutputProperties
		if _, ok := props.Datasource.AsDocumentDbOutputDataSource(); !ok {
			return fmt.Errorf("specified output is not of type %s", streamanalytics.TypeBasicOutputDataSourceTypeMicrosoftStorageDocumentDB)
		}
		return nil
	}
}
