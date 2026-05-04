package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlMISyncTaskOutput = MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel{}

type MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel struct {
	ActiveBackupSets          *[]BackupSetInfo        `json:"activeBackupSets,omitempty"`
	ContainerName             *string                 `json:"containerName,omitempty"`
	EndedOn                   *string                 `json:"endedOn,omitempty"`
	ErrorPrefix               *string                 `json:"errorPrefix,omitempty"`
	ExceptionsAndWarnings     *[]ReportableException  `json:"exceptionsAndWarnings,omitempty"`
	FullBackupSetInfo         *BackupSetInfo          `json:"fullBackupSetInfo,omitempty"`
	IsFullBackupRestored      *bool                   `json:"isFullBackupRestored,omitempty"`
	LastRestoredBackupSetInfo *BackupSetInfo          `json:"lastRestoredBackupSetInfo,omitempty"`
	MigrationState            *DatabaseMigrationState `json:"migrationState,omitempty"`
	SourceDatabaseName        *string                 `json:"sourceDatabaseName,omitempty"`
	StartedOn                 *string                 `json:"startedOn,omitempty"`

	// Fields inherited from MigrateSqlServerSqlMISyncTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel) MigrateSqlServerSqlMISyncTaskOutput() BaseMigrateSqlServerSqlMISyncTaskOutputImpl {
	return BaseMigrateSqlServerSqlMISyncTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel{}

func (s MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlMISyncTaskOutputDatabaseLevel: %+v", err)
	}

	return encoded, nil
}
