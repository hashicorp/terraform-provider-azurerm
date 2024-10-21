package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbTaskOutput = MigrateSqlServerSqlDbTaskOutputMigrationLevel{}

type MigrateSqlServerSqlDbTaskOutputMigrationLevel struct {
	DatabaseSummary           *map[string]DatabaseSummaryResult `json:"databaseSummary,omitempty"`
	Databases                 *map[string]string                `json:"databases,omitempty"`
	DurationInSeconds         *int64                            `json:"durationInSeconds,omitempty"`
	EndedOn                   *string                           `json:"endedOn,omitempty"`
	ExceptionsAndWarnings     *[]ReportableException            `json:"exceptionsAndWarnings,omitempty"`
	Message                   *string                           `json:"message,omitempty"`
	MigrationReportResult     *MigrationReportResult            `json:"migrationReportResult,omitempty"`
	MigrationValidationResult *MigrationValidationResult        `json:"migrationValidationResult,omitempty"`
	SourceServerBrandVersion  *string                           `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion       *string                           `json:"sourceServerVersion,omitempty"`
	StartedOn                 *string                           `json:"startedOn,omitempty"`
	Status                    *MigrationStatus                  `json:"status,omitempty"`
	StatusMessage             *string                           `json:"statusMessage,omitempty"`
	TargetServerBrandVersion  *string                           `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion       *string                           `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbTaskOutputMigrationLevel) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskOutputMigrationLevel{}

func (s MigrateSqlServerSqlDbTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
