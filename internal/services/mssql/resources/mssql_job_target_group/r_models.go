// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_target_group

type MsSqlJobTargetGroupResourceModel struct {
	Name       string           `tfschema:"name"`
	JobAgentID string           `tfschema:"job_agent_id"`
	JobTargets []MsSqlJobTarget `tfschema:"job_target"`
}

func (r Resource) ModelObject() any {
	return &MsSqlJobTargetGroupResourceModel{}
}
