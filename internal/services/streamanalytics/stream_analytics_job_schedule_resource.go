// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	streamAnalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type JobScheduleResource struct{}

type JobScheduleResourceModel struct {
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_id"`
	StartMode          string `tfschema:"start_mode"`
	StartTime          string `tfschema:"start_time"`
	LastOutputTime     string `tfschema:"last_output_time"`
}

var _ sdk.ResourceWithStateMigration = JobScheduleResource{}

func (r JobScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: streamingjobs.ValidateStreamingJobID,
		},

		"start_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(streamingjobs.OutputStartModeCustomTime),
				string(streamingjobs.OutputStartModeJobStartTime),
				string(streamingjobs.OutputStartModeLastOutputEventTime),
			}, false),
		},

		"start_time": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ISO8601DateTime,
		},
	}
}

func (r JobScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_output_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r JobScheduleResource) ModelObject() interface{} {
	return &JobScheduleResourceModel{}
}

func (r JobScheduleResource) ResourceType() string {
	return "azurerm_stream_analytics_job_schedule"
}

func (r JobScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return streamAnalyticsValidate.StreamingJobScheduleID
}

func (r JobScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model JobScheduleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.JobsClient
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
					return fmt.Errorf("`start_mode` can only be set to `LastOutputEventTime` if this job was previously started")
				}
			}

			props := &streamingjobs.StartStreamingJobParameters{
				OutputStartMode: utils.ToPtr(outputStartMode),
			}

			if outputStartMode == streamingjobs.OutputStartModeCustomTime {
				if model.StartTime == "" {
					return fmt.Errorf("`start_time` must be specified if `start_mode` is set to `CustomTime`")
				} else {
					startTime, _ := date.ParseTime(time.RFC3339, model.StartTime)
					outputStartTime := &date.Time{
						Time: startTime,
					}
					props.OutputStartTime = utils.String(outputStartTime.String())
				}
			}

			if err := client.StartThenPoll(ctx, *streamAnalyticsId, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r JobScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient
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
					startTime := ""
					if v := props.OutputStartTime; v != nil {
						startTime = *v
					}

					lastOutputTime := ""
					if v := props.LastOutputEventTime; v != nil {
						lastOutputTime = *v
					}

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

func (r JobScheduleResource) Update() sdk.ResourceFunc {
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
					OutputStartMode: utils.ToPtr(outputStartMode),
				}

				if outputStartMode == streamingjobs.OutputStartModeCustomTime {
					props.OutputStartTime = utils.String(outputStartTime.String())
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

func (r JobScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient
			id, err := parse.StreamingJobScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingJobName)
			if err := client.StopThenPoll(ctx, streamingJobId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r JobScheduleResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsJobScheduleV0ToV1{},
		},
	}
}
