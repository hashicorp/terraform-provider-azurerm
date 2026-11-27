// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_step

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobsteps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Elastic Job Step.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"job_id": {
			Description:  "The ID of the Elastic Job.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: jobsteps.ValidateJobID,
		},
		"job_credential_id": {
			Description:  "The ID of the Elastic Job Credential to use when executing this Elastic Job Step. Omit this argument to run the step under the Job Agent's Managed Identity (user-assigned).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: jobcredentials.ValidateCredentialID,
		},
		"job_step_index": {
			Description:  "The index at which to insert this Elastic Job Step into the Elastic Job.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"job_target_group_id": {
			Description:  "The ID of the Elastic Job Target Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: jobtargetgroups.ValidateTargetGroupID,
		},
		"sql_script": {
			Description:  "The T-SQL script to be executed by this Elastic Job Step.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"initial_retry_interval_seconds": {
			Description:  "The initial retry interval in seconds.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 2147483),
		},
		"maximum_retry_interval_seconds": {
			Description:  "The maximum retry interval in seconds.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      120,
			ValidateFunc: validation.IntBetween(1, 2147483),
		},
		"output_target": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mssql_database_id": {
						Description:  "The ID of the output database.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSqlDatabaseID,
					},
					"table_name": {
						Description:  "The name of the output table.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"job_credential_id": {
						Description:  "The ID of the Elastic Job Credential to use when connecting to the output destination.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: jobcredentials.ValidateCredentialID,
					},
					"schema_name": {
						Description:  "The name of the output schema.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "dbo",
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"retry_attempts": {
			Description:  "The number of retry attempts.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(1, math.MaxInt32),
		},
		"retry_interval_backoff_multiplier": {
			Description:  "The multiplier for time between retries.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      2.0,
			ValidateFunc: validation.FloatAtLeast(1),
		},
		"timeout_seconds": {
			Description:  "The execution timeout in seconds for this Elastic Job Step.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      43200,
			ValidateFunc: validation.IntBetween(1, 2147483),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlJobStepResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.MaximumRetryIntervalSeconds <= model.InitialRetryIntervalSeconds {
				return errors.New("`maximum_retry_interval_seconds` must be greater than `initial_retry_interval_seconds`")
			}

			// Once set, `job_credential_id` cannot be removed
			// https://github.com/Azure/azure-rest-api-specs/issues/35881
			if o, n := metadata.ResourceDiff.GetChange("job_credential_id"); o.(string) != "" && n.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("job_credential_id"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
