package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImportSettings = SnowflakeImportCopyCommand{}

type SnowflakeImportCopyCommand struct {
	AdditionalCopyOptions   *map[string]interface{} `json:"additionalCopyOptions,omitempty"`
	AdditionalFormatOptions *map[string]interface{} `json:"additionalFormatOptions,omitempty"`
	StorageIntegration      *string                 `json:"storageIntegration,omitempty"`

	// Fields inherited from ImportSettings

	Type string `json:"type"`
}

func (s SnowflakeImportCopyCommand) ImportSettings() BaseImportSettingsImpl {
	return BaseImportSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SnowflakeImportCopyCommand{}

func (s SnowflakeImportCopyCommand) MarshalJSON() ([]byte, error) {
	type wrapper SnowflakeImportCopyCommand
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SnowflakeImportCopyCommand: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SnowflakeImportCopyCommand: %+v", err)
	}

	decoded["type"] = "SnowflakeImportCopyCommand"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SnowflakeImportCopyCommand: %+v", err)
	}

	return encoded, nil
}
