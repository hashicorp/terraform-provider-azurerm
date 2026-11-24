// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule

type MsSqlJobScheduleResourceModel struct {
	JobID     string `tfschema:"job_id"`
	Type      string `tfschema:"type"`
	Enabled   bool   `tfschema:"enabled"`
	EndTime   string `tfschema:"end_time"`
	Interval  string `tfschema:"interval"`
	StartTime string `tfschema:"start_time"`
}

func (Resource) ModelObject() any {
	return &MsSqlJobScheduleResourceModel{}
}
