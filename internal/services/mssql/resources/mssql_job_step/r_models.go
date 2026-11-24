// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_step

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

func (r Resource) ModelObject() any {
	return &MsSqlJobStepResourceModel{}
}
