package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportSettings interface {
	ImportSettings() BaseImportSettingsImpl
}

var _ ImportSettings = BaseImportSettingsImpl{}

type BaseImportSettingsImpl struct {
	Type string `json:"type"`
}

func (s BaseImportSettingsImpl) ImportSettings() BaseImportSettingsImpl {
	return s
}

var _ ImportSettings = RawImportSettingsImpl{}

// RawImportSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawImportSettingsImpl struct {
	importSettings BaseImportSettingsImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawImportSettingsImpl) ImportSettings() BaseImportSettingsImpl {
	return s.importSettings
}

func UnmarshalImportSettingsImplementation(input []byte) (ImportSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ImportSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLakeImportCommand") {
		var out AzureDatabricksDeltaLakeImportCommand
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeImportCommand: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeImportCopyCommand") {
		var out SnowflakeImportCopyCommand
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeImportCopyCommand: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TeradataImportCommand") {
		var out TeradataImportCommand
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeradataImportCommand: %+v", err)
		}
		return out, nil
	}

	var parent BaseImportSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseImportSettingsImpl: %+v", err)
	}

	return RawImportSettingsImpl{
		importSettings: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
