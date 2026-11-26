// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/startstopmanagedinstanceschedules"
)

type SqlManagedInstanceStartStopScheduleModel struct {
	SqlManagedInstanceId string              `tfschema:"managed_instance_id"`
	Description          string              `tfschema:"description"`
	Schedule             []ScheduleItemModel `tfschema:"schedule"`
	TimeZoneId           string              `tfschema:"timezone_id"`
	NextExecutionTime    string              `tfschema:"next_execution_time"`
	NextRunAction        string              `tfschema:"next_run_action"`
}

type ScheduleItemModel struct {
	StartDay  startstopmanagedinstanceschedules.DayOfWeek `tfschema:"start_day"`
	StartTime string                                      `tfschema:"start_time"`
	StopDay   startstopmanagedinstanceschedules.DayOfWeek `tfschema:"stop_day"`
	StopTime  string                                      `tfschema:"stop_time"`
}

func (r Resource) ModelObject() any {
	return &SqlManagedInstanceStartStopScheduleModel{}
}
