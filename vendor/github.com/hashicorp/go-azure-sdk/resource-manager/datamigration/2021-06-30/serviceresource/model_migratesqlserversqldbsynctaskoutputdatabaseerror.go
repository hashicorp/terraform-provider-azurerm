package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbSyncTaskOutput = MigrateSqlServerSqlDbSyncTaskOutputDatabaseError{}

type MigrateSqlServerSqlDbSyncTaskOutputDatabaseError struct {
	ErrorMessage *string                            `json:"errorMessage,omitempty"`
	Events       *[]SyncMigrationDatabaseErrorEvent `json:"events,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbSyncTaskOutputDatabaseError) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbSyncTaskOutputDatabaseError{}

func (s MigrateSqlServerSqlDbSyncTaskOutputDatabaseError) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbSyncTaskOutputDatabaseError
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseError: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseError: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelErrorOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseError: %+v", err)
	}

	return encoded, nil
}
