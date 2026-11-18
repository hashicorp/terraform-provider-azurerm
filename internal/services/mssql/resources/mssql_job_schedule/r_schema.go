// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validation.ISO8601Duration,
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

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
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
