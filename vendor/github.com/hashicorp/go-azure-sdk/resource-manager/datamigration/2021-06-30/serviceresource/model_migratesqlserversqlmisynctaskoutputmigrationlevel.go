package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMISyncTaskOutput = MigrateSqlServerSqlMISyncTaskOutputMigrationLevel{}

type MigrateSqlServerSqlMISyncTaskOutputMigrationLevel struct {
	DatabaseCount            *int64          `json:"databaseCount,omitempty"`
	DatabaseErrorCount       *int64          `json:"databaseErrorCount,omitempty"`
	EndedOn                  *string         `json:"endedOn,omitempty"`
	SourceServerBrandVersion *string         `json:"sourceServerBrandVersion,omitempty"`
	SourceServerName         *string         `json:"sourceServerName,omitempty"`
	SourceServerVersion      *string         `json:"sourceServerVersion,omitempty"`
	StartedOn                *string         `json:"startedOn,omitempty"`
	State                    *MigrationState `json:"state,omitempty"`
	TargetServerBrandVersion *string         `json:"targetServerBrandVersion,omitempty"`
	TargetServerName         *string         `json:"targetServerName,omitempty"`
	TargetServerVersion      *string         `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMISyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMISyncTaskOutputMigrationLevel) MigrateSqlServerSqlMISyncTaskOutput() BaseMigrateSqlServerSqlMISyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlMISyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMISyncTaskOutputMigrationLevel{}

func (s MigrateSqlServerSqlMISyncTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMISyncTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMISyncTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMISyncTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMISyncTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
