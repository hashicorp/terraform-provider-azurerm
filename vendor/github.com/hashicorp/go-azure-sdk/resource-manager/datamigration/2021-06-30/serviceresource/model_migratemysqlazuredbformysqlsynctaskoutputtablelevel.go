package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateMySqlAzureDbForMySqlSyncTaskOutput = MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel{}

type MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel struct {
	CdcDeleteCounter      *string                  `json:"cdcDeleteCounter,omitempty"`
	CdcInsertCounter      *string                  `json:"cdcInsertCounter,omitempty"`
	CdcUpdateCounter      *string                  `json:"cdcUpdateCounter,omitempty"`
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

	// Fields inherited from MigrateMySqlAzureDbForMySqlSyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel) MigrateMySqlAzureDbForMySqlSyncTaskOutput() BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl {
	return BaseMigrateMySqlAzureDbForMySqlSyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel{}

func (s MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel: %+v", err)
	}

	decoded["resultType"] = "TableLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateMySqlAzureDbForMySqlSyncTaskOutputTableLevel: %+v", err)
	}

	return encoded, nil
}
