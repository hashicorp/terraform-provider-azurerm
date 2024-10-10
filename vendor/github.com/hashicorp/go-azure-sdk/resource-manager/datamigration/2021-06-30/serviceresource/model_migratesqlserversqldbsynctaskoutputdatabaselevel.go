package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbSyncTaskOutput = MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel{}

type MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel struct {
	AppliedChanges          *int64                               `json:"appliedChanges,omitempty"`
	CdcDeleteCounter        *int64                               `json:"cdcDeleteCounter,omitempty"`
	CdcInsertCounter        *int64                               `json:"cdcInsertCounter,omitempty"`
	CdcUpdateCounter        *int64                               `json:"cdcUpdateCounter,omitempty"`
	DatabaseName            *string                              `json:"databaseName,omitempty"`
	EndedOn                 *string                              `json:"endedOn,omitempty"`
	FullLoadCompletedTables *int64                               `json:"fullLoadCompletedTables,omitempty"`
	FullLoadErroredTables   *int64                               `json:"fullLoadErroredTables,omitempty"`
	FullLoadLoadingTables   *int64                               `json:"fullLoadLoadingTables,omitempty"`
	FullLoadQueuedTables    *int64                               `json:"fullLoadQueuedTables,omitempty"`
	IncomingChanges         *int64                               `json:"incomingChanges,omitempty"`
	InitializationCompleted *bool                                `json:"initializationCompleted,omitempty"`
	Latency                 *int64                               `json:"latency,omitempty"`
	MigrationState          *SyncDatabaseMigrationReportingState `json:"migrationState,omitempty"`
	StartedOn               *string                              `json:"startedOn,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel{}

func (s MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbSyncTaskOutputDatabaseLevel: %+v", err)
	}

	return encoded, nil
}
