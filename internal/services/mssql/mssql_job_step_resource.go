// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobsteps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlJobStepResource struct{}

type MsSqlJobStepResourceModel struct {
	Name                           string                `tfschema:"name"`
	JobID                          string                `tfschema:"job_id"`
	JobCredentialID                string                `tfschema:"job_credential_id"`
	JobStepIndex                   int64                 `tfschema:"job_step_index"`
	JobTargetGroupID               string                `tfschema:"job_target_group_id"`
	SqlScript                      string                `tfschema:"sql_script"`
	InitialRetryIntervalSeconds    int64                 `tfschema:"initial_retry_interval_seconds"`
	MaximumRetryIntervalSeconds    int64                 `tfschema:"maximum_retry_interval_seconds"`
	OutputTarget                   []JobStepOutputTarget `tfschema:"output_target"`
	RetryAttempts                  int64                 `tfschema:"retry_attempts"`
	RetryIntervalBackoffMultiplier float64               `tfschema:"retry_interval_backoff_multiplier"`
	TimeoutSeconds                 int64                 `tfschema:"timeout_seconds"`
}

type JobStepOutputTarget struct {
	JobCredentialId string `tfschema:"job_credential_id"`
	MsSqlDatabaseId string `tfschema:"mssql_database_id"`
	TableName       string `tfschema:"table_name"`
	SchemaName      string `tfschema:"schema_name"`
}

var (
	_ sdk.ResourceWithUpdate        = MsSqlJobStepResource{}
	_ sdk.ResourceWithCustomizeDiff = MsSqlJobStepResource{}
)

func (MsSqlJobStepResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},
		"job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: jobsteps.ValidateJobID,
			ForceNew:     true,
		},
		"job_credential_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: jobcredentials.ValidateCredentialID,
		},
		"job_step_index": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"job_target_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: jobtargetgroups.ValidateTargetGroupID,
		},
		"sql_script": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"initial_retry_interval_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 2147483),
		},
		"maximum_retry_interval_seconds": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSqlDatabaseID,
					},
					"table_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"job_credential_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: jobcredentials.ValidateCredentialID,
					},
					"schema_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "dbo",
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"retry_attempts": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(1, math.MaxInt32),
		},
		"retry_interval_backoff_multiplier": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      2.0,
			ValidateFunc: validation.FloatAtLeast(1),
		},
		"timeout_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      43200,
			ValidateFunc: validation.IntBetween(1, 2147483),
		},
	}
}

func (r MsSqlJobStepResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlJobStepResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlJobStepResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.MaximumRetryIntervalSeconds <= model.InitialRetryIntervalSeconds {
				return fmt.Errorf("`maximum_retry_interval_seconds` must be greater than `initial_retry_interval_seconds`")
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

func (r MsSqlJobStepResource) ModelObject() interface{} {
	return &MsSqlJobStepResourceModel{}
}

func (r MsSqlJobStepResource) ResourceType() string {
	return "azurerm_mssql_job_step"
}

func (r MsSqlJobStepResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

func (r MsSqlJobStepResource) Read() sdk.ResourceFunc {
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

					state.JobStepIndex = pointer.From(props.StepId)
					state.JobTargetGroupID = props.TargetGroup
					state.SqlScript = props.Action.Value
					state.InitialRetryIntervalSeconds = pointer.From(props.ExecutionOptions.InitialRetryIntervalSeconds)
					state.MaximumRetryIntervalSeconds = pointer.From(props.ExecutionOptions.MaximumRetryIntervalSeconds)

					target, err := flattenOutputTarget(props.Output)
					if err != nil {
						return fmt.Errorf("flattening `output_target`: %+v", err)
					}
					state.OutputTarget = target

					state.RetryAttempts = pointer.From(props.ExecutionOptions.RetryAttempts)
					state.RetryIntervalBackoffMultiplier = pointer.From(props.ExecutionOptions.RetryIntervalBackoffMultiplier)
					state.TimeoutSeconds = pointer.From(props.ExecutionOptions.TimeoutSeconds)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MsSqlJobStepResource) Update() sdk.ResourceFunc {
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

func (r MsSqlJobStepResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobStepsClient

			id, err := jobsteps.ParseStepID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlJobStepResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobsteps.ValidateStepID
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

func stringPtrIfSet(s string) *string {
	if s == "" {
		return nil
	}
	return pointer.To(s)
}
