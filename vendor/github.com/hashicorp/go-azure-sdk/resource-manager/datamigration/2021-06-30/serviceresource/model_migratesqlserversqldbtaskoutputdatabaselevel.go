package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbTaskOutput = MigrateSqlServerSqlDbTaskOutputDatabaseLevel{}

type MigrateSqlServerSqlDbTaskOutputDatabaseLevel struct {
	DatabaseName             *string                                    `json:"databaseName,omitempty"`
	EndedOn                  *string                                    `json:"endedOn,omitempty"`
	ErrorCount               *int64                                     `json:"errorCount,omitempty"`
	ErrorPrefix              *string                                    `json:"errorPrefix,omitempty"`
	ExceptionsAndWarnings    *[]ReportableException                     `json:"exceptionsAndWarnings,omitempty"`
	Message                  *string                                    `json:"message,omitempty"`
	NumberOfObjects          *int64                                     `json:"numberOfObjects,omitempty"`
	NumberOfObjectsCompleted *int64                                     `json:"numberOfObjectsCompleted,omitempty"`
	ObjectSummary            *map[string]DataItemMigrationSummaryResult `json:"objectSummary,omitempty"`
	ResultPrefix             *string                                    `json:"resultPrefix,omitempty"`
	Stage                    *DatabaseMigrationStage                    `json:"stage,omitempty"`
	StartedOn                *string                                    `json:"startedOn,omitempty"`
	State                    *MigrationState                            `json:"state,omitempty"`
	StatusMessage            *string                                    `json:"statusMessage,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbTaskOutputDatabaseLevel) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskOutputDatabaseLevel{}

func (s MigrateSqlServerSqlDbTaskOutputDatabaseLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskOutputDatabaseLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevel: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevel: %+v", err)
	}

	return encoded, nil
}
