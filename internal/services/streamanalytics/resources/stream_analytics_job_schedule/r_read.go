// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient
			id, err := parse.StreamingJobScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			streamAnalyticsId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingJobName)

			var opts streamingjobs.GetOperationOptions
			resp, err := client.Get(ctx, streamAnalyticsId, opts)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					startTime := pointer.From(props.OutputStartTime)

					lastOutputTime := pointer.From(props.LastOutputEventTime)

					startMode := ""
					if v := props.OutputStartMode; v != nil {
						startMode = string(*v)
					}

					state := JobScheduleResourceModel{
						StreamAnalyticsJob: streamAnalyticsId.ID(),
						StartMode:          startMode,
						StartTime:          startTime,
						LastOutputTime:     lastOutputTime,
					}

					return metadata.Encode(&state)
				}
			}
			return nil
		},
	}
}
