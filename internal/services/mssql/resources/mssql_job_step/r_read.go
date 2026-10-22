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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobsteps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobStepsClient

			id, err := jobsteps.ParseStepID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := MsSqlJobStepResourceModel{
				Name:  id.StepName,
				JobID: jobsteps.NewJobID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName, id.JobName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if v := pointer.From(props.Credential); v != "" {
						credentialID, err := jobcredentials.ParseCredentialID(v)
						if err != nil {
							return err
						}
						state.JobCredentialID = credentialID.ID()
					}

					target, err := flattenOutputTarget(props.Output)
					if err != nil {
						return fmt.Errorf("flattening `output_target`: %+v", err)
					}
					state.OutputTarget = target

					state.JobStepIndex = pointer.From(props.StepId)
					state.JobTargetGroupID = props.TargetGroup
					state.SqlScript = props.Action.Value

					if exec := props.ExecutionOptions; exec != nil {
						state.InitialRetryIntervalSeconds = pointer.From(exec.InitialRetryIntervalSeconds)
						state.MaximumRetryIntervalSeconds = pointer.From(exec.MaximumRetryIntervalSeconds)
						state.RetryAttempts = pointer.From(exec.RetryAttempts)
						state.RetryIntervalBackoffMultiplier = pointer.From(exec.RetryIntervalBackoffMultiplier)
						state.TimeoutSeconds = pointer.From(exec.TimeoutSeconds)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenOutputTarget(input *jobsteps.JobStepOutput) ([]JobStepOutputTarget, error) {
	if input == nil {
		return []JobStepOutputTarget{}, nil
	}

	credentialID := ""
	if v := pointer.From(input.Credential); v != "" {
		id, err := jobcredentials.ParseCredentialID(v)
		if err != nil {
			return nil, err
		}
		credentialID = id.ID()
	}

	databaseId := commonids.NewSqlDatabaseID(pointer.From(input.SubscriptionId), pointer.From(input.ResourceGroupName), input.ServerName, input.DatabaseName)

	return []JobStepOutputTarget{
		{
			JobCredentialId: credentialID,
			MsSqlDatabaseId: databaseId.ID(),
			TableName:       input.TableName,
			SchemaName:      pointer.From(input.SchemaName),
		},
	}, nil
}

type JobStepOutputTarget struct {
	JobCredentialId string `tfschema:"job_credential_id"`
	MsSqlDatabaseId string `tfschema:"mssql_database_id"`
	TableName       string `tfschema:"table_name"`
	SchemaName      string `tfschema:"schema_name"`
}

func stringPtrIfSet(s string) *string {
	if s == "" {
		return nil
	}
	return pointer.To(s)
}
