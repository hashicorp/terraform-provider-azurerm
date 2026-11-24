// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

type NetAppSnapshotPolicyResourceModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	AccountName       string                 `tfschema:"account_name"`
	Enabled           bool                   `tfschema:"enabled"`
	HourlySchedule    []HourlyScheduleModel  `tfschema:"hourly_schedule"`
	DailySchedule     []DailyScheduleModel   `tfschema:"daily_schedule"`
	WeeklySchedule    []WeeklyScheduleModel  `tfschema:"weekly_schedule"`
	MonthlySchedule   []MonthlyScheduleModel `tfschema:"monthly_schedule"`
	Tags              map[string]string      `tfschema:"tags"`
}

type HourlyScheduleModel struct {
	SnapshotsToKeep int64 `tfschema:"snapshots_to_keep"`
	Minute          int64 `tfschema:"minute"`
}

type DailyScheduleModel struct {
	SnapshotsToKeep int64 `tfschema:"snapshots_to_keep"`
	Hour            int64 `tfschema:"hour"`
	Minute          int64 `tfschema:"minute"`
}

type WeeklyScheduleModel struct {
	SnapshotsToKeep int64    `tfschema:"snapshots_to_keep"`
	DaysOfWeek      []string `tfschema:"days_of_week"`
	Hour            int64    `tfschema:"hour"`
	Minute          int64    `tfschema:"minute"`
}

type MonthlyScheduleModel struct {
	SnapshotsToKeep int64   `tfschema:"snapshots_to_keep"`
	DaysOfMonth     []int64 `tfschema:"days_of_month"`
	Hour            int64   `tfschema:"hour"`
	Minute          int64   `tfschema:"minute"`
}

func (Resource) ModelObject() any {
	return &NetAppSnapshotPolicyResourceModel{}
}
