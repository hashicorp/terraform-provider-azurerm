// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient
			id, err := parse.StreamingJobScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state JobScheduleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChanges("start_mode", "start_time") {
				outputStartMode := streamingjobs.OutputStartMode(state.StartMode)
				startTime, _ := date.ParseTime(time.RFC3339, state.StartTime)
				outputStartTime := &date.Time{
					Time: startTime,
				}

				props := &streamingjobs.StartStreamingJobParameters{
					OutputStartMode: pointer.To(outputStartMode),
				}

				if outputStartMode == streamingjobs.OutputStartModeCustomTime {
					props.OutputStartTime = pointer.To(outputStartTime.String())
				}

				var opts streamingjobs.GetOperationOptions
				streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingJobName)
				existing, err := client.Get(ctx, streamingJobId, opts)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}

				if v := existing.Model.Properties; v != nil && v.JobState != nil && *v.JobState == "Running" {
					if err := client.StopThenPoll(ctx, streamingJobId); err != nil {
						return fmt.Errorf("stopping %s: %+v", *id, err)
					}
				}

				if err := client.StartThenPoll(ctx, streamingJobId, *props); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
