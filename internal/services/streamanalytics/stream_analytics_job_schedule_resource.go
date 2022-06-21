package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2020-03-01/streamanalytics"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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

func (r JobScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: streamAnalyticsValidate.StreamingJobID,
		},

		"start_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(streamanalytics.OutputStartModeCustomTime),
				string(streamanalytics.OutputStartModeJobStartTime),
				string(streamanalytics.OutputStartModeLastOutputEventTime),
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
			streamAnalyticsId, err := parse.StreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}

			// This is a virtual resource so the last segment is hardcoded
			id := parse.NewStreamingJobScheduleID(streamAnalyticsId.SubscriptionId, streamAnalyticsId.ResourceGroup, streamAnalyticsId.Name, "default")

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, "")
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			outputStartMode := streamanalytics.OutputStartMode(model.StartMode)
			if outputStartMode == streamanalytics.OutputStartModeLastOutputEventTime {
				if v := existing.StreamingJobProperties.LastOutputEventTime; v == nil {
					return fmt.Errorf("`start_mode` can only be set to `LastOutputEventTime` if this job was previously started")
				}
			}

			props := &streamanalytics.StartStreamingJobParameters{
				OutputStartMode: outputStartMode,
			}

			if outputStartMode == streamanalytics.OutputStartModeCustomTime {
				if model.StartTime == "" {
					return fmt.Errorf("`start_time` must be specified if `start_mode` is set to `CustomTime`")
				} else {
					startTime, _ := date.ParseTime(time.RFC3339, model.StartTime)
					outputStartTime := &date.Time{
						Time: startTime,
					}
					props.OutputStartTime = outputStartTime
				}
			}

			future, err := client.Start(ctx, id.ResourceGroup, id.StreamingjobName, props)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting on create/update for %s: %+v", id, err)
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

			streamAnalyticsId := parse.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingjobName)

			resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if props := resp.StreamingJobProperties; props != nil {
				startTime := ""
				if v := props.OutputStartTime; v != nil {
					startTime = v.String()
				}

				lastOutputTime := ""
				if v := props.LastOutputEventTime; v != nil {
					lastOutputTime = v.String()
				}

				state := JobScheduleResourceModel{
					StreamAnalyticsJob: streamAnalyticsId.ID(),
					StartMode:          string(props.OutputStartMode),
					StartTime:          startTime,
					LastOutputTime:     lastOutputTime,
				}

				return metadata.Encode(&state)
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
				outputStartMode := streamanalytics.OutputStartMode(state.StartMode)
				startTime, _ := date.ParseTime(time.RFC3339, state.StartTime)
				outputStartTime := &date.Time{
					Time: startTime,
				}

				props := &streamanalytics.StartStreamingJobParameters{
					OutputStartMode: outputStartMode,
				}

				if outputStartMode == streamanalytics.OutputStartModeCustomTime {
					props.OutputStartTime = outputStartTime
				}

				existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, "")
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}

				if v := existing.StreamingJobProperties; v != nil && v.JobState != nil && *v.JobState == "Running" {
					future, err := client.Stop(ctx, id.ResourceGroup, id.StreamingjobName)
					if err != nil {
						return err
					}
					if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for %s to stop: %+v", *id, err)
					}
				}

				future, err := client.Start(ctx, id.ResourceGroup, id.StreamingjobName, props)
				if err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
				if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for update of %q: %+v", *id, err)
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

			future, err := client.Stop(ctx, id.ResourceGroup, id.StreamingjobName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}
			return nil
		},
	}
}
