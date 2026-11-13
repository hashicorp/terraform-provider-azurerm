// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_step

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobsteps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobStepsClient

			var model MsSqlJobStepResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			job, err := jobsteps.ParseJobID(model.JobID)
			if err != nil {
				return err
			}

			id := jobsteps.NewStepID(job.SubscriptionId, job.ResourceGroupName, job.ServerName, job.JobAgentName, job.JobName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := jobsteps.JobStep{
				Name: pointer.To(model.Name),
				Properties: pointer.To(jobsteps.JobStepProperties{
					Action: jobsteps.JobStepAction{
						Value: model.SqlScript,
					},
					Credential: stringPtrIfSet(model.JobCredentialID),
					ExecutionOptions: pointer.To(jobsteps.JobStepExecutionOptions{
						InitialRetryIntervalSeconds:    pointer.To(model.InitialRetryIntervalSeconds),
						MaximumRetryIntervalSeconds:    pointer.To(model.MaximumRetryIntervalSeconds),
						RetryAttempts:                  pointer.To(model.RetryAttempts),
						RetryIntervalBackoffMultiplier: pointer.To(model.RetryIntervalBackoffMultiplier),
						TimeoutSeconds:                 pointer.To(model.TimeoutSeconds),
					}),
					StepId:      pointer.To(model.JobStepIndex),
					TargetGroup: model.JobTargetGroupID,
				}),
			}

			target, err := expandOutputTarget(model.OutputTarget)
			if err != nil {
				return fmt.Errorf("expanding `output_target`: %+v", err)
			}
			parameters.Properties.Output = target

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandOutputTarget(input []JobStepOutputTarget) (*jobsteps.JobStepOutput, error) {
	if len(input) == 0 {
		return nil, nil
	}

	target := input[0]
	databaseId, err := commonids.ParseSqlDatabaseID(target.MsSqlDatabaseId)
	if err != nil {
		return nil, err
	}

	return pointer.To(jobsteps.JobStepOutput{
		Credential:        pointer.To(target.JobCredentialId),
		DatabaseName:      databaseId.DatabaseName,
		ResourceGroupName: pointer.To(databaseId.ResourceGroupName),
		SchemaName:        pointer.To(target.SchemaName),
		ServerName:        databaseId.ServerName,
		SubscriptionId:    pointer.To(databaseId.SubscriptionId),
		TableName:         target.TableName,
	}), nil
}
