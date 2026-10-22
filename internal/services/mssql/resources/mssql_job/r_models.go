// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job

type MsSqlJobResourceModel struct {
	Name        string `tfschema:"name"`
	JobAgentID  string `tfschema:"job_agent_id"`
	Description string `tfschema:"description"`
}

func (Resource) ModelObject() any {
	return &MsSqlJobResourceModel{}
}
