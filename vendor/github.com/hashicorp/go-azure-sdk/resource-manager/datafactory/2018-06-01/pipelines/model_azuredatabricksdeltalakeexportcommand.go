package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ExportSettings = AzureDatabricksDeltaLakeExportCommand{}

type AzureDatabricksDeltaLakeExportCommand struct {
	DateFormat      *interface{} `json:"dateFormat,omitempty"`
	TimestampFormat *interface{} `json:"timestampFormat,omitempty"`

	// Fields inherited from ExportSettings

	Type string `json:"type"`
}

func (s AzureDatabricksDeltaLakeExportCommand) ExportSettings() BaseExportSettingsImpl {
	return BaseExportSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureDatabricksDeltaLakeExportCommand{}

func (s AzureDatabricksDeltaLakeExportCommand) MarshalJSON() ([]byte, error) {
	type wrapper AzureDatabricksDeltaLakeExportCommand
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDatabricksDeltaLakeExportCommand: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDatabricksDeltaLakeExportCommand: %+v", err)
	}

	decoded["type"] = "AzureDatabricksDeltaLakeExportCommand"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDatabricksDeltaLakeExportCommand: %+v", err)
	}

	return encoded, nil
}
