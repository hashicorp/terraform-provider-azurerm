package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CommandProperties = MigrateSyncCompleteCommandProperties{}

type MigrateSyncCompleteCommandProperties struct {
	Input  *MigrateSyncCompleteCommandInput  `json:"input,omitempty"`
	Output *MigrateSyncCompleteCommandOutput `json:"output,omitempty"`

	// Fields inherited from CommandProperties

	CommandType string        `json:"commandType"`
	Errors      *[]ODataError `json:"errors,omitempty"`
	State       *CommandState `json:"state,omitempty"`
}

func (s MigrateSyncCompleteCommandProperties) CommandProperties() BaseCommandPropertiesImpl {
	return BaseCommandPropertiesImpl{
		CommandType: s.CommandType,
		Errors:      s.Errors,
		State:       s.State,
	}
}

var _ json.Marshaler = MigrateSyncCompleteCommandProperties{}

func (s MigrateSyncCompleteCommandProperties) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSyncCompleteCommandProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSyncCompleteCommandProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSyncCompleteCommandProperties: %+v", err)
	}

	decoded["commandType"] = "Migrate.Sync.Complete.Database"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSyncCompleteCommandProperties: %+v", err)
	}

	return encoded, nil
}
