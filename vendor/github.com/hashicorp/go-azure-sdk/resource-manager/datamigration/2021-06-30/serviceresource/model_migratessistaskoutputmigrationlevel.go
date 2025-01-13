package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSsisTaskOutput = MigrateSsisTaskOutputMigrationLevel{}

type MigrateSsisTaskOutputMigrationLevel struct {
	EndedOn                  *string                `json:"endedOn,omitempty"`
	ExceptionsAndWarnings    *[]ReportableException `json:"exceptionsAndWarnings,omitempty"`
	Message                  *string                `json:"message,omitempty"`
	SourceServerBrandVersion *string                `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion      *string                `json:"sourceServerVersion,omitempty"`
	Stage                    *SsisMigrationStage    `json:"stage,omitempty"`
	StartedOn                *string                `json:"startedOn,omitempty"`
	Status                   *MigrationStatus       `json:"status,omitempty"`
	TargetServerBrandVersion *string                `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion      *string                `json:"targetServerVersion,omitempty"`

	// Fields inherited from MigrateSsisTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSsisTaskOutputMigrationLevel) MigrateSsisTaskOutput() BaseMigrateSsisTaskOutputImpl {
	return BaseMigrateSsisTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSsisTaskOutputMigrationLevel{}

func (s MigrateSsisTaskOutputMigrationLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSsisTaskOutputMigrationLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSsisTaskOutputMigrationLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSsisTaskOutputMigrationLevel: %+v", err)
	}

	decoded["resultType"] = "MigrationLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSsisTaskOutputMigrationLevel: %+v", err)
	}

	return encoded, nil
}
