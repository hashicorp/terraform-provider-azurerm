package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbTaskOutput = MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult{}

type MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult struct {
	DataIntegrityValidationResult *DataIntegrityValidationResult    `json:"dataIntegrityValidationResult,omitempty"`
	EndedOn                       *string                           `json:"endedOn,omitempty"`
	MigrationId                   *string                           `json:"migrationId,omitempty"`
	QueryAnalysisValidationResult *QueryAnalysisValidationResult    `json:"queryAnalysisValidationResult,omitempty"`
	SchemaValidationResult        *SchemaComparisonValidationResult `json:"schemaValidationResult,omitempty"`
	SourceDatabaseName            *string                           `json:"sourceDatabaseName,omitempty"`
	StartedOn                     *string                           `json:"startedOn,omitempty"`
	Status                        *ValidationStatus                 `json:"status,omitempty"`
	TargetDatabaseName            *string                           `json:"targetDatabaseName,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult{}

func (s MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult: %+v", err)
	}

	decoded["resultType"] = "MigrationDatabaseLevelValidationOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskOutputDatabaseLevelValidationResult: %+v", err)
	}

	return encoded, nil
}
