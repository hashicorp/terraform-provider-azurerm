package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectToSourceSqlServerTaskOutput = ConnectToSourceSqlServerTaskOutputDatabaseLevel{}

type ConnectToSourceSqlServerTaskOutputDatabaseLevel struct {
	CompatibilityLevel *DatabaseCompatLevel `json:"compatibilityLevel,omitempty"`
	DatabaseFiles      *[]DatabaseFileInfo  `json:"databaseFiles,omitempty"`
	DatabaseState      *DatabaseState       `json:"databaseState,omitempty"`
	Name               *string              `json:"name,omitempty"`
	SizeMB             *float64             `json:"sizeMB,omitempty"`

	// Fields inherited from ConnectToSourceSqlServerTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s ConnectToSourceSqlServerTaskOutputDatabaseLevel) ConnectToSourceSqlServerTaskOutput() BaseConnectToSourceSqlServerTaskOutputImpl {
	return BaseConnectToSourceSqlServerTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = ConnectToSourceSqlServerTaskOutputDatabaseLevel{}

func (s ConnectToSourceSqlServerTaskOutputDatabaseLevel) MarshalJSON() ([]byte, error) {
	type wrapper ConnectToSourceSqlServerTaskOutputDatabaseLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConnectToSourceSqlServerTaskOutputDatabaseLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectToSourceSqlServerTaskOutputDatabaseLevel: %+v", err)
	}

	decoded["resultType"] = "DatabaseLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConnectToSourceSqlServerTaskOutputDatabaseLevel: %+v", err)
	}

	return encoded, nil
}
