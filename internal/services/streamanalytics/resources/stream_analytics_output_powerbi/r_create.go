// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_powerbi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
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
			var model OutputPowerBIResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := streamanalyticsClient.From(metadata.Client).OutputsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			streamingJobId, err := streamingjobs.ParseStreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}
			id := outputs.NewOutputID(subscriptionId, streamingJobId.ResourceGroupName, streamingJobId.StreamingJobName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			powerBIOutputProps := &outputs.PowerBIOutputDataSourceProperties{
				Dataset:            pointer.To(model.DataSet),
				Table:              pointer.To(model.Table),
				GroupId:            pointer.To(model.GroupID),
				GroupName:          pointer.To(model.GroupName),
				RefreshToken:       pointer.To("someRefreshToken"),                    // A valid refresh token is currently only obtainable via the Azure Portal. Put a dummy string value here when creating the data source and then going to the Azure Portal to authenticate the data source which will update this property with a valid refresh token.
				AuthenticationMode: pointer.ToEnum[outputs.AuthenticationMode]("Msi"), // Set authentication mode as "Msi" here since other modes requires params obtainable from portal only.
			}

			if model.TokenUserDisplayName != "" {
				powerBIOutputProps.TokenUserDisplayName = pointer.To(model.TokenUserDisplayName)
			}

			if model.TokenUserPrincipalName != "" {
				powerBIOutputProps.TokenUserPrincipalName = pointer.To(model.TokenUserPrincipalName)
			}

			props := outputs.Output{
				Name: pointer.To(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.PowerBIOutputDataSource{
						Properties: powerBIOutputProps,
					},
				},
			}

			var opts outputs.CreateOrReplaceOperationOptions
			if _, err = client.CreateOrReplace(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
