// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
					output, ok := props.Datasource.(outputs.AzureTableOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a Table Output", *id)
					}

					if output.Properties != nil {
						if output.Properties.AccountName == nil || output.Properties.Table == nil || output.Properties.PartitionKey == nil || output.Properties.RowKey == nil || output.Properties.BatchSize == nil {
							return nil
						}

						state := OutputTableResourceModel{
							Name:               id.OutputName,
							ResourceGroup:      id.ResourceGroupName,
							StreamAnalyticsJob: id.StreamingJobName,
							StorageAccountKey:  metadata.ResourceData.Get("storage_account_key").(string),
						}

						state.StorageAccount = pointer.From(output.Properties.AccountName)

						state.Table = pointer.From(output.Properties.Table)

						state.PartitionKey = pointer.From(output.Properties.PartitionKey)

						state.RowKey = pointer.From(output.Properties.RowKey)

						state.BatchSize = pointer.From(output.Properties.BatchSize)

						var columnsToRemove []string
						if columns := output.Properties.ColumnsToRemove; columns != nil && len(*columns) > 0 {
							columnsToRemove = *columns
						}
						state.ColumnsToRemove = columnsToRemove

						return metadata.Encode(&state)
					}
				}
			}
			return nil
		},
	}
}
