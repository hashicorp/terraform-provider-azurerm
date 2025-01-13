package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateOracleAzureDbPostgreSqlSyncTaskOutput = MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError{}

type MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError struct {
	ErrorMessage *string                            `json:"errorMessage,omitempty"`
	Events       *[]SyncMigrationDatabaseErrorEvent `json:"events,omitempty"`

	// Fields inherited from MigrateOracleAzureDbPostgreSqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError) MigrateOracleAzureDbPostgreSqlSyncTaskOutput() BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl {
	return BaseMigrateOracleAzureDbPostgreSqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError{}

func (s MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateOracleAzureDbPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	return encoded, nil
}
