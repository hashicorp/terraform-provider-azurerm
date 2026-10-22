// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_database

type MsSqlManagedDatabaseModel struct {
	Name                    string                    `tfschema:"name"`
	ManagedInstanceId       string                    `tfschema:"managed_instance_id"`
	LongTermRetentionPolicy []LongTermRetentionPolicy `tfschema:"long_term_retention_policy"`
	ShortTermRetentionDays  int64                     `tfschema:"short_term_retention_days"`
	PointInTimeRestore      []PointInTimeRestore      `tfschema:"point_in_time_restore"`
	Tags                    map[string]string         `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MsSqlManagedDatabaseModel{}
}
