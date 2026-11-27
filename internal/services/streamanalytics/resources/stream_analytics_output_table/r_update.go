// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

			var state OutputTableResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := &outputs.AzureTableOutputDataSourceProperties{
				AccountName:  pointer.To(state.StorageAccount),
				AccountKey:   pointer.To(state.StorageAccountKey),
				Table:        pointer.To(state.Table),
				PartitionKey: pointer.To(state.PartitionKey),
				RowKey:       pointer.To(state.RowKey),
				BatchSize:    pointer.To(state.BatchSize),
			}

			if metadata.ResourceData.HasChange("columns_to_remove") {
				props.ColumnsToRemove = &state.ColumnsToRemove
			}

			output := outputs.Output{
				Name: pointer.To(state.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureTableOutputDataSource{
						Properties: props,
					},
				},
			}
			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, output, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
