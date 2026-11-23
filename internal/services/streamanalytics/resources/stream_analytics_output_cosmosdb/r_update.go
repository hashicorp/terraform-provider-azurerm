// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_cosmosdb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			databaseId, err := cosmosdb.ParseSqlDatabaseID(state.Database)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChangesExcept("name", "stream_analytics_job_id") {
				props := outputs.Output{
					Properties: &outputs.OutputProperties{
						Datasource: outputs.DocumentDbOutputDataSource{
							Properties: &outputs.DocumentDbOutputDataSourceProperties{
								AccountKey:            &state.AccountKey,
								Database:              &databaseId.SqlDatabaseName,
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
