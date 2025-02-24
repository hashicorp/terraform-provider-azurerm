package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ExportSettings = SnowflakeExportCopyCommand{}

type SnowflakeExportCopyCommand struct {
	AdditionalCopyOptions   *map[string]interface{} `json:"additionalCopyOptions,omitempty"`
	AdditionalFormatOptions *map[string]interface{} `json:"additionalFormatOptions,omitempty"`
	StorageIntegration      *string                 `json:"storageIntegration,omitempty"`

	// Fields inherited from ExportSettings

	Type string `json:"type"`
}

func (s SnowflakeExportCopyCommand) ExportSettings() BaseExportSettingsImpl {
	return BaseExportSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SnowflakeExportCopyCommand{}

func (s SnowflakeExportCopyCommand) MarshalJSON() ([]byte, error) {
	type wrapper SnowflakeExportCopyCommand
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SnowflakeExportCopyCommand: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SnowflakeExportCopyCommand: %+v", err)
	}

	decoded["type"] = "SnowflakeExportCopyCommand"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SnowflakeExportCopyCommand: %+v", err)
	}

	return encoded, nil
}
