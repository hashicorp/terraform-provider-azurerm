// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cosmosParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	cosmosValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OutputCosmosDBResource struct{}

var (
	_ sdk.ResourceWithCustomImporter = OutputCosmosDBResource{}
	_ sdk.ResourceWithStateMigration = OutputCosmosDBResource{}
)

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
			ValidateFunc: streamingjobs.ValidateStreamingJobID,
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

			streamingJobId, err := streamingjobs.ParseStreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}
			id := outputs.NewOutputID(subscriptionId, streamingJobId.ResourceGroupName, streamingJobId.StreamingJobName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			databaseId, err := cosmosParse.SqlDatabaseID(model.Database)
			if err != nil {
				return err
			}

			documentDbOutputProps := &outputs.DocumentDbOutputDataSourceProperties{
				AccountId:             utils.String(databaseId.DatabaseAccountName),
				AccountKey:            utils.String(model.AccountKey),
				Database:              utils.String(databaseId.Name),
				CollectionNamePattern: utils.String(model.ContainerName),
				DocumentId:            utils.String(model.DocumentID),
				PartitionKey:          utils.String(model.PartitionKey),
			}

			props := outputs.Output{
				Name: utils.String(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.DocumentDbOutputDataSource{
						Properties: documentDbOutputProps,
					},
				},
			}

			var opts outputs.CreateOrReplaceOperationOptions
			if _, err := client.CreateOrReplace(ctx, id, props, opts); err != nil {
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

			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					output, ok := props.Datasource.(outputs.DocumentDbOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a CosmosDb Output", *id)
					}

					streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)
					state := OutputCosmosDBResourceModel{
						Name:               id.OutputName,
						StreamAnalyticsJob: streamingJobId.ID(),
					}

					state.AccountKey = metadata.ResourceData.Get("cosmosdb_account_key").(string)

					databaseId := cosmosParse.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroupName, *output.Properties.AccountId, *output.Properties.Database)
					state.Database = databaseId.ID()

					collectionName := ""
					if v := output.Properties.CollectionNamePattern; v != nil {
						collectionName = *v
					}
					state.ContainerName = collectionName

					document := ""
					if v := output.Properties.DocumentId; v != nil {
						document = *v
					}
					state.DocumentID = document

					partitionKey := ""
					if v := output.Properties.PartitionKey; v != nil {
						partitionKey = *v
					}
					state.PartitionKey = partitionKey

					return metadata.Encode(&state)
				}
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
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r OutputCosmosDBResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outputs.ValidateOutputID
}

func (r OutputCosmosDBResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
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
				props := outputs.Output{
					Properties: &outputs.OutputProperties{
						Datasource: outputs.DocumentDbOutputDataSource{
							Properties: &outputs.DocumentDbOutputDataSourceProperties{
								AccountKey:            &state.AccountKey,
								Database:              &databaseId.Name,
								CollectionNamePattern: &state.ContainerName,
								DocumentId:            &state.DocumentID,
								PartitionKey:          &state.PartitionKey,
							},
						},
					},
				}
				var opts outputs.UpdateOperationOptions
				if _, err := client.Update(ctx, *id, props, opts); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r OutputCosmosDBResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.Model.Properties
		if _, ok := props.Datasource.(outputs.DocumentDbOutputDataSource); !ok {
			return fmt.Errorf("specified output is not of type")
		}
		return nil
	}
}

func (r OutputCosmosDBResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputCosmosDbV0ToV1{},
		},
	}
}
