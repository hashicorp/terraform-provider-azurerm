package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbSyncTaskOutput = MigrateSqlServerSqlDbSyncTaskOutputTableLevel{}

type MigrateSqlServerSqlDbSyncTaskOutputTableLevel struct {
	CdcDeleteCounter      *int64                   `json:"cdcDeleteCounter,omitempty"`
	CdcInsertCounter      *int64                   `json:"cdcInsertCounter,omitempty"`
	CdcUpdateCounter      *int64                   `json:"cdcUpdateCounter,omitempty"`
	DataErrorsCounter     *int64                   `json:"dataErrorsCounter,omitempty"`
	DatabaseName          *string                  `json:"databaseName,omitempty"`
	FullLoadEndedOn       *string                  `json:"fullLoadEndedOn,omitempty"`
	FullLoadEstFinishTime *string                  `json:"fullLoadEstFinishTime,omitempty"`
	FullLoadStartedOn     *string                  `json:"fullLoadStartedOn,omitempty"`
	FullLoadTotalRows     *int64                   `json:"fullLoadTotalRows,omitempty"`
	LastModifiedTime      *string                  `json:"lastModifiedTime,omitempty"`
	State                 *SyncTableMigrationState `json:"state,omitempty"`
	TableName             *string                  `json:"tableName,omitempty"`
	TotalChangesApplied   *int64                   `json:"totalChangesApplied,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbSyncTaskOutputTableLevel) MigrateSqlServerSqlDbSyncTaskOutput() BaseMigrateSqlServerSqlDbSyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbSyncTaskOutputTableLevel{}

func (s MigrateSqlServerSqlDbSyncTaskOutputTableLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbSyncTaskOutputTableLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbSyncTaskOutputTableLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbSyncTaskOutputTableLevel: %+v", err)
	}

	decoded["resultType"] = "TableLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbSyncTaskOutputTableLevel: %+v", err)
	}

	return encoded, nil
}
