package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationValidationDatabaseSummaryResult struct {
	EndedOn            *string           `json:"endedOn,omitempty"`
	Id                 *string           `json:"id,omitempty"`
	MigrationId        *string           `json:"migrationId,omitempty"`
	SourceDatabaseName *string           `json:"sourceDatabaseName,omitempty"`
	StartedOn          *string           `json:"startedOn,omitempty"`
	Status             *ValidationStatus `json:"status,omitempty"`
	TargetDatabaseName *string           `json:"targetDatabaseName,omitempty"`
}

func (o *MigrationValidationDatabaseSummaryResult) GetEndedOnAsTime() (*time.Time, error) {
	if o.EndedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrationValidationDatabaseSummaryResult) SetEndedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndedOn = &formatted
}

func (o *MigrationValidationDatabaseSummaryResult) GetStartedOnAsTime() (*time.Time, error) {
	if o.StartedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrationValidationDatabaseSummaryResult) SetStartedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedOn = &formatted
}
