// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_function

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
					output, ok := props.Datasource.(outputs.AzureFunctionOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a Function Output", *id)
					}

					if output.Properties != nil {
						if output.Properties.FunctionAppName == nil || output.Properties.FunctionName == nil || output.Properties.MaxBatchCount == nil || output.Properties.MaxBatchSize == nil {
							return nil
						}

						state := OutputFunctionResourceModel{
							Name:               id.OutputName,
							ResourceGroup:      id.ResourceGroupName,
							StreamAnalyticsJob: id.StreamingJobName,
							ApiKey:             metadata.ResourceData.Get("api_key").(string),
							FunctionApp:        pointer.From(output.Properties.FunctionAppName),
							FunctionName:       pointer.From(output.Properties.FunctionName),
							BatchMaxInBytes:    int64(pointer.From(output.Properties.MaxBatchSize)),
							BatchMaxCount:      int64(pointer.From(output.Properties.MaxBatchCount)),
						}

						return metadata.Encode(&state)
					}
				}
			}
			return nil
		},
	}
}
