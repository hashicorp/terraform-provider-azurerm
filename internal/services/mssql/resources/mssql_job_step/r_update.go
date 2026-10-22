// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_step

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobsteps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobStepsClient

			id, err := jobsteps.ParseStepID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config MsSqlJobStepResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
			}
			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("job_credential_id") {
				props.Credential = stringPtrIfSet(config.JobCredentialID)
			}

			if metadata.ResourceData.HasChange("job_step_index") {
				props.StepId = pointer.To(config.JobStepIndex)
			}

			if metadata.ResourceData.HasChange("job_target_group_id") {
				props.TargetGroup = config.JobTargetGroupID
			}

			if metadata.ResourceData.HasChange("sql_script") {
				props.Action.Value = config.SqlScript
			}

			if metadata.ResourceData.HasChange("initial_retry_interval_seconds") {
				props.ExecutionOptions.InitialRetryIntervalSeconds = pointer.To(config.InitialRetryIntervalSeconds)
			}

			if metadata.ResourceData.HasChange("maximum_retry_interval_seconds") {
				props.ExecutionOptions.MaximumRetryIntervalSeconds = pointer.To(config.MaximumRetryIntervalSeconds)
			}

			if metadata.ResourceData.HasChange("output_target") {
				target, err := expandOutputTarget(config.OutputTarget)
				if err != nil {
					return fmt.Errorf("expanding `output_target`: %+v", err)
				}
				props.Output = target
			}

			if metadata.ResourceData.HasChange("retry_attempts") {
				props.ExecutionOptions.RetryAttempts = pointer.To(config.RetryAttempts)
			}

			if metadata.ResourceData.HasChange("retry_interval_backoff_multiplier") {
				props.ExecutionOptions.RetryIntervalBackoffMultiplier = pointer.To(config.RetryIntervalBackoffMultiplier)
			}

			if metadata.ResourceData.HasChange("timeout_seconds") {
				props.ExecutionOptions.TimeoutSeconds = pointer.To(config.TimeoutSeconds)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
