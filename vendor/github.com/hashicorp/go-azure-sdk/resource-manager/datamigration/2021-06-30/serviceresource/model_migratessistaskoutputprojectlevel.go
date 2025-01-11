package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSsisTaskOutput = MigrateSsisTaskOutputProjectLevel{}

type MigrateSsisTaskOutputProjectLevel struct {
	EndedOn               *string                `json:"endedOn,omitempty"`
	ExceptionsAndWarnings *[]ReportableException `json:"exceptionsAndWarnings,omitempty"`
	FolderName            *string                `json:"folderName,omitempty"`
	Message               *string                `json:"message,omitempty"`
	ProjectName           *string                `json:"projectName,omitempty"`
	Stage                 *SsisMigrationStage    `json:"stage,omitempty"`
	StartedOn             *string                `json:"startedOn,omitempty"`
	State                 *MigrationState        `json:"state,omitempty"`

	// Fields inherited from MigrateSsisTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSsisTaskOutputProjectLevel) MigrateSsisTaskOutput() BaseMigrateSsisTaskOutputImpl {
	return BaseMigrateSsisTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSsisTaskOutputProjectLevel{}

func (s MigrateSsisTaskOutputProjectLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSsisTaskOutputProjectLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSsisTaskOutputProjectLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSsisTaskOutputProjectLevel: %+v", err)
	}

	decoded["resultType"] = "SsisProjectLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSsisTaskOutputProjectLevel: %+v", err)
	}

	return encoded, nil
}
