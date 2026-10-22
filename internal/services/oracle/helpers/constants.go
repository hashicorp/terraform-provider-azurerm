// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

type LongTermBackUpScheduleDetails struct {
	RepeatCadence         string `tfschema:"repeat_cadence"`
	TimeOfBackup          string `tfschema:"time_of_backup"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`
	Enabled               bool   `tfschema:"enabled"`
}
