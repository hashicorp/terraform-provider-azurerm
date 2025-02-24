// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlJobScheduleResource struct{}

type MsSqlJobScheduleResourceModel struct {
	JobID     string `tfschema:"job_id"`
	Type      string `tfschema:"type"`
	Enabled   bool   `tfschema:"enabled"`
	EndTime   string `tfschema:"end_time"`
	Interval  string `tfschema:"interval"`
	StartTime string `tfschema:"start_time"`
}

var (
	_ sdk.ResourceWithUpdate        = MsSqlJobScheduleResource{}
	_ sdk.ResourceWithCustomizeDiff = MsSqlJobScheduleResource{}
)

func (MsSqlJobScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: jobs.ValidateJobID,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(jobs.PossibleValuesForJobScheduleType(), false),
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// Note: O+C the API sets this field to `false` if the schedule type is `Once` and the job has finished
			Computed: true,
		},
		"end_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C API sets and returns a default value if omitted
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},
		"interval": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ISO8601Duration,
		},
		"start_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C API sets and returns a default value if omitted
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},
	}
}

func (MsSqlJobScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (MsSqlJobScheduleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config MsSqlJobScheduleResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if config.Type == string(jobs.JobScheduleTypeRecurring) && config.Interval == "" {
				return fmt.Errorf("`interval` must be set when `type` is `Recurring`")
			}

			return nil
		},
	}
}

func (MsSqlJobScheduleResource) ModelObject() interface{} {
	return &MsSqlJobScheduleResourceModel{}
}

func (MsSqlJobScheduleResource) ResourceType() string {
	return "azurerm_mssql_job_schedule"
}

func (r MsSqlJobScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			var config MsSqlJobScheduleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			jobId, err := jobs.ParseJobID(config.JobID)
			if err != nil {
				return err
			}

			locks.ByID(jobId.ID())
			defer locks.UnlockByID(jobId.ID())

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found: %+v", jobId, err)
				}

				return fmt.Errorf("checking for presence of existing %s: %+v", jobId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", jobId)
			}

			if existing.Model.Properties.Schedule == nil {
				return fmt.Errorf("retrieving %s: `model.Properties.Schedule` was nil", jobId)
			}

			// Default schedule is disabled when created using the API
			// if schedule is enabled we can reasonably assume the schedule was modified outside of Terraform and should be imported.
			schedule := existing.Model.Properties.Schedule
			if pointer.From(schedule.Enabled) {
				return metadata.ResourceRequiresImport(r.ResourceType(), jobId)
			}

			schedule.Enabled = pointer.To(config.Enabled)
			schedule.Type = pointer.To(jobs.JobScheduleType(config.Type))

			if config.EndTime != "" {
				schedule.EndTime = pointer.To(config.EndTime)
			}
			if config.Interval != "" && config.Type == string(jobs.JobScheduleTypeRecurring) {
				schedule.Interval = pointer.To(config.Interval)
			}
			if config.StartTime != "" {
				schedule.StartTime = pointer.To(config.StartTime)
			}

			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("creating schedule for %s: %+v", jobId, err)
			}

			metadata.SetID(jobId)
			return nil
		},
	}
}

func (MsSqlJobScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			jobId, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(jobId)
				}

				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			state := MsSqlJobScheduleResourceModel{
				JobID: jobId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if schedule := props.Schedule; schedule != nil {
						state.Enabled = pointer.From(schedule.Enabled)
						state.EndTime = pointer.From(schedule.EndTime)
						state.Interval = pointer.From(schedule.Interval)
						state.StartTime = pointer.From(schedule.StartTime)
						state.Type = string(pointer.From(schedule.Type))
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (MsSqlJobScheduleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			jobId, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(jobId.ID())
			defer locks.UnlockByID(jobId.ID())

			var config MsSqlJobScheduleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", jobId)
			}

			if existing.Model.Properties.Schedule == nil {
				return fmt.Errorf("retrieving %s: `schedule` was nil", jobId)
			}

			schedule := existing.Model.Properties.Schedule
			if metadata.ResourceData.HasChange("enabled") {
				schedule.Enabled = pointer.To(config.Enabled)
			}

			if metadata.ResourceData.HasChange("end_time") {
				schedule.EndTime = pointer.To(config.EndTime)
			}

			if metadata.ResourceData.HasChange("interval") {
				schedule.Interval = pointer.To(config.Interval)
			}

			if metadata.ResourceData.HasChange("start_time") {
				schedule.StartTime = pointer.To(config.StartTime)
			}

			if metadata.ResourceData.HasChange("type") {
				schedule.Type = pointer.To(jobs.JobScheduleType(config.Type))
			}

			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("updating schedule for %s: %+v", jobId, err)
			}

			return nil
		},
	}
}

func (MsSqlJobScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			jobId, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(jobId.ID())
			defer locks.UnlockByID(jobId.ID())

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return nil
				}

				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			if model := existing.Model; model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", jobId)
			}

			// Set Schedule to nil allowing Azure to repopulate with default values
			existing.Model.Properties.Schedule = nil
			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("deleting schedule from %s, %+v", jobId, err)
			}

			return nil
		},
	}
}

func (MsSqlJobScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobs.ValidateJobID
}
