// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model JobScheduleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := streamanalyticsClient.From(metadata.Client).JobsClient
			streamAnalyticsId, err := streamingjobs.ParseStreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}

			// This is a virtual resource so the last segment is hardcoded
			id := parse.NewStreamingJobScheduleID(streamAnalyticsId.SubscriptionId, streamAnalyticsId.ResourceGroupName, streamAnalyticsId.StreamingJobName, "default")

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			var opts streamingjobs.GetOperationOptions
			existing, err := client.Get(ctx, *streamAnalyticsId, opts)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			outputStartMode := streamingjobs.OutputStartMode(model.StartMode)
			if outputStartMode == streamingjobs.OutputStartModeLastOutputEventTime {
				if v := existing.Model.Properties.LastOutputEventTime; v == nil {
					return errors.New("`start_mode` can only be set to `LastOutputEventTime` if this job was previously started")
				}
			}

			props := &streamingjobs.StartStreamingJobParameters{
				OutputStartMode: pointer.To(outputStartMode),
			}

			if outputStartMode == streamingjobs.OutputStartModeCustomTime {
				if model.StartTime == "" {
					return errors.New("`start_time` must be specified if `start_mode` is set to `CustomTime`")
				} else {
					startTime, _ := date.ParseTime(time.RFC3339, model.StartTime)
					outputStartTime := &date.Time{
						Time: startTime,
					}
					props.OutputStartTime = pointer.To(outputStartTime.String())
				}
			}

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			if err := client.StartThenPoll(ctx, *streamAnalyticsId, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
