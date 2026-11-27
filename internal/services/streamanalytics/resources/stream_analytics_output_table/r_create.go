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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := streamanalyticsClient.From(metadata.Client).OutputsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			id := outputs.NewOutputID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsJob, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			tableOutputProps := &outputs.AzureTableOutputDataSourceProperties{
				AccountName:  pointer.To(model.StorageAccount),
				AccountKey:   pointer.To(model.StorageAccountKey),
				Table:        pointer.To(model.Table),
				PartitionKey: pointer.To(model.PartitionKey),
				RowKey:       pointer.To(model.RowKey),
				BatchSize:    pointer.To(model.BatchSize),
			}

			if v := model.ColumnsToRemove; len(v) > 0 {
				tableOutputProps.ColumnsToRemove = &v
			}

			props := outputs.Output{
				Name: pointer.To(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.AzureTableOutputDataSource{
						Properties: tableOutputProps,
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
