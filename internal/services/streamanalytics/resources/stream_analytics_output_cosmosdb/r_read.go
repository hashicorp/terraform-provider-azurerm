// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_cosmosdb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).OutputsClient

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

					databaseId := cosmosdb.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroupName, *output.Properties.AccountId, *output.Properties.Database)
					state.Database = databaseId.ID()

					state.ContainerName = pointer.From(output.Properties.CollectionNamePattern)

					state.DocumentID = pointer.From(output.Properties.DocumentId)

					state.PartitionKey = pointer.From(output.Properties.PartitionKey)

					state.AuthenticationMode = pointer.FromEnum(output.Properties.AuthenticationMode)

					return metadata.Encode(&state)
				}
			}
			return nil
		},
	}
}
