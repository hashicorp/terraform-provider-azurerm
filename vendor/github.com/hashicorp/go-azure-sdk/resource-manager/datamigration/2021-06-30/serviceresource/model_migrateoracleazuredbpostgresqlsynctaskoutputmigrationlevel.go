package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateOracleAzureDbPostgreSqlSyncTaskOutput = MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel{}

type MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel struct {
	EndedOn             *string `json:"endedOn,omitempty"`
	SourceServer        *string `json:"sourceServer,omitempty"`
	SourceServerVersion *string `json:"sourceServerVersion,omitempty"`
	StartedOn           *string `json:"startedOn,omitempty"`
	TargetServer        *string `json:"targetServer,omitempty"`
	TargetServerVersion *string `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateOracleAzureDbPostgreSqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel) MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl {
	return BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel{}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
