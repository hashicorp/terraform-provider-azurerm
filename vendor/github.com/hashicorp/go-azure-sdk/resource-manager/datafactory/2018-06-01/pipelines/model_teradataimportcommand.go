package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImportSettings = TeradataImportCommand{}

type TeradataImportCommand struct {
	AdditionalFormatOptions *map[string]string `json:"additionalFormatOptions,omitempty"`

	// Fields inherited from ImportSettings

	Type string `json:"type"`
}

func (s TeradataImportCommand) ImportSettings() BaseImportSettingsImpl {
	return BaseImportSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = TeradataImportCommand{}

func (s TeradataImportCommand) MarshalJSON() ([]byte, error) {
	type wrapper TeradataImportCommand
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TeradataImportCommand: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TeradataImportCommand: %+v", err)
	}

	decoded["type"] = "TeradataImportCommand"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TeradataImportCommand: %+v", err)
	}

	return encoded, nil
}
