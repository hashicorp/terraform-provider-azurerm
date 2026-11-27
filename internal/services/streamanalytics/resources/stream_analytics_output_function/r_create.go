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
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputFunctionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.OutputsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := outputs.NewOutputID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsJob, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := outputs.Output{
				Name: pointer.To(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureFunctionOutputDataSource{
						Properties: &outputs.AzureFunctionOutputDataSourceProperties{
							FunctionAppName: pointer.To(model.FunctionApp),
							FunctionName:    pointer.To(model.FunctionName),
							ApiKey:          pointer.To(model.ApiKey),
							MaxBatchSize:    pointer.To(float64(model.BatchMaxInBytes)),
							MaxBatchCount:   pointer.To(float64(model.BatchMaxCount)),
						},
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
