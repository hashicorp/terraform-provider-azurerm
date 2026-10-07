package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput = MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError{}

type MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError struct {
	ErrorMessage *string                            `json:"errorMessage,omitempty"`
	Events       *[]SyncMigrationDatabaseErrorEvent `json:"events,omitempty"`

	// Fields inherited from MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError) MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutput() BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl {
	return BaseMigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError{}

func (s MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError) MarshalJSON() ([]byte, error) {
	type wrapper MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigratePostgreSqlAzureDbForPostgreSqlSyncTaskOutputDatabaseError: %+v", err)
	}

	return encoded, nil
}
