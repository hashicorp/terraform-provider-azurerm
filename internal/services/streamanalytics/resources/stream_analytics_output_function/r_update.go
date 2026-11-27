// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_function

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputFunctionResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := outputs.Output{
				Name: pointer.To(state.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureFunctionOutputDataSource{
						Properties: &outputs.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: pointer.To(state.FunctionApp),
							FunctionName:    pointer.To(state.FunctionName),
							ApiKey:          pointer.To(state.ApiKey),
							MaxBatchSize:    pointer.To(float64(state.BatchMaxInBytes)),
							MaxBatchCount:   pointer.To(float64(state.BatchMaxCount)),
						},
					},
				},
			}

			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, props, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
