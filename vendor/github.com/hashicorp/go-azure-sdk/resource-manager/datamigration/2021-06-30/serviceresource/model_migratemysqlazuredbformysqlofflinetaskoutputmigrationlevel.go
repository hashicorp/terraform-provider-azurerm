package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateMySqlAzureDbForMySqlOfflineTaskOutput = MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel{}

type MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel struct {
	DatabaseSummary          *map[string]DatabaseSummaryResult `json:"databaseSummary,omitempty"`
	Databases                *map[string]string                `json:"databases,omitempty"`
	DurationInSeconds        *int64                            `json:"durationInSeconds,omitempty"`
	EndedOn                  *string                           `json:"endedOn,omitempty"`
	ExceptionsAndWarnings    *[]ReportableException            `json:"exceptionsAndWarnings,omitempty"`
	LastStorageUpdate        *string                           `json:"lastStorageUpdate,omitempty"`
	Message                  *string                           `json:"message,omitempty"`
	MigrationReportResult    *MigrationReportResult            `json:"migrationReportResult,omitempty"`
	SourceServerBrandVersion *string                           `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion      *string                           `json:"sourceServerVersion,omitempty"`
	StartedOn                *string                           `json:"startedOn,omitempty"`
	Status                   *MigrationStatus                  `json:"status,omitempty"`
	StatusMessage            *string                           `json:"statusMessage,omitempty"`
	TargetServerBrandVersion *string                           `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion      *string                           `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateMySqlAzureDbForMySqlOfflineTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel) MigrateMySqlAzureDbForMySqlOfflineTaskOutput() BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl {
	return BaseMigrateMySqlAzureDbForMySqlOfflineTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel{}

func (s MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateMySqlAzureDbForMySqlOfflineTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
