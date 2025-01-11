package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbSyncTaskOutput = MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel{}

type MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel struct {
	DatabaseCount       *int64  `json:"databaseCount,omitempty"`
	EndedOn             *string `json:"endedOn,omitempty"`
	SourceServer        *string `json:"sourceServer,omitempty"`
	SourceServerVersion *string `json:"sourceServerVersion,omitempty"`
	StartedOn           *string `json:"startedOn,omitempty"`
	TargetServer        *string `json:"targetServer,omitempty"`
	TargetServerVersion *string `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel{}

func (s MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbSyncTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
