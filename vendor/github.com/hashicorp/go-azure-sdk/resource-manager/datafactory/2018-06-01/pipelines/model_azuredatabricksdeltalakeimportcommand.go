package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImportSettings = AzureDatabricksDeltaLakeImportCommand{}

type AzureDatabricksDeltaLakeImportCommand struct {
	DateFormat      *string `json:"dateFormat,omitempty"`
	TimestampFormat *string `json:"timestampFormat,omitempty"`

	// Fields inherited from ImportSettings

	Type string `json:"type"`
}

func (s AzureDatabricksDeltaLakeImportCommand) ImportSettings() BaseImportSettingsImpl {
	return BaseImportSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureDatabricksDeltaLakeImportCommand{}

func (s AzureDatabricksDeltaLakeImportCommand) MarshalJSON() ([]byte, error) {
	type wrapper AzureDatabricksDeltaLakeImportCommand
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDatabricksDeltaLakeImportCommand: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDatabricksDeltaLakeImportCommand: %+v", err)
	}

	decoded["type"] = "AzureDatabricksDeltaLakeImportCommand"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDatabricksDeltaLakeImportCommand: %+v", err)
	}

	return encoded, nil
}
