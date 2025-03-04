package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMITaskOutput = MigrateSqlServerSqlMITaskOutputDatabaseLevel{}

type MigrateSqlServerSqlMITaskOutputDatabaseLevel struct {
	DatabaseName          *string                 `json:"databaseName,omitempty"`
	EndedOn               *string                 `json:"endedOn,omitempty"`
	ExceptionsAndWarnings *[]ReportableException  `json:"exceptionsAndWarnings,omitempty"`
	Message               *string                 `json:"message,omitempty"`
	SizeMB                *float64                `json:"sizeMB,omitempty"`
	Stage                 *DatabaseMigrationStage `json:"stage,omitempty"`
	StartedOn             *string                 `json:"startedOn,omitempty"`
	State                 *MigrationState         `json:"state,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMITaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMITaskOutputDatabaseLevel) MigrateSqlServerSqlMITaskOutput() BaseMigrateSqlServerSqlMITaskOutputImpl {
	return BaseMigrateSqlServerSqlMITaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMITaskOutputDatabaseLevel{}

func (s MigrateSqlServerSqlMITaskOutputDatabaseLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMITaskOutputDatabaseLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMITaskOutputDatabaseLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMITaskOutputDatabaseLevel: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMITaskOutputDatabaseLevel: %+v", err)
	}

	return encoded, nil
}
