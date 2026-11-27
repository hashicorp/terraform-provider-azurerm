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
					output, ok := props.Datasource.(outputs.PowerBIOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a PowerBI Output", *id)
					}

					streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)

					state := OutputPowerBIResourceModel{
						Name:               id.OutputName,
						StreamAnalyticsJob: streamingJobId.ID(),
					}

					state.DataSet = pointer.From(output.Properties.Dataset)

					state.Table = pointer.From(output.Properties.Table)

					state.GroupID = pointer.From(output.Properties.GroupId)

					state.GroupName = pointer.From(output.Properties.GroupName)

					state.TokenUserDisplayName = metadata.ResourceData.Get("token_user_display_name").(string)
					state.TokenUserPrincipalName = metadata.ResourceData.Get("token_user_principal_name").(string)

					return metadata.Encode(&state)
				}
			}
			return nil
		},
	}
}
