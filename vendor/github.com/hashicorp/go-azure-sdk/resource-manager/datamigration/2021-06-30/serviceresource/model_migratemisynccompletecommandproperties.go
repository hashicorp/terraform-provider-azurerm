package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CommandProperties = MigrateMISyncCompleteCommandProperties{}

type MigrateMISyncCompleteCommandProperties struct {
	Input  *MigrateMISyncCompleteCommandInput  `json:"input,omitempty"`
	Output *MigrateMISyncCompleteCommandOutput `json:"output,omitempty"`

	// Fields inherited from CommandProperties

	CommandType string        `json:"commandType"`
	Errors      *[]ODataError `json:"errors,omitempty"`
	State       *CommandState `json:"state,omitempty"`
}

func (s MigrateMISyncCompleteCommandProperties) CommandProperties() BaseCommandPropertiesImpl {
	return BaseCommandPropertiesImpl{
		CommandType: s.CommandType,
		Errors:      s.Errors,
		State:       s.State,
	}
}

var _ json.Marshaler = MigrateMISyncCompleteCommandProperties{}

func (s MigrateMISyncCompleteCommandProperties) MarshalJSON() ([]byte, error) {
	type wrapper MigrateMISyncCompleteCommandProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateMISyncCompleteCommandProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateMISyncCompleteCommandProperties: %+v", err)
	}

	decoded["commandType"] = "Migrate.SqlServer.AzureDbSqlMi.Complete"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateMISyncCompleteCommandProperties: %+v", err)
	}

	return encoded, nil
}
