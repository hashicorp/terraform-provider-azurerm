package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportSettings interface {
	ExportSettings() BaseExportSettingsImpl
}

var _ ExportSettings = BaseExportSettingsImpl{}

type BaseExportSettingsImpl struct {
	Type string `json:"type"`
}

func (s BaseExportSettingsImpl) ExportSettings() BaseExportSettingsImpl {
	return s
}

var _ ExportSettings = RawExportSettingsImpl{}

// RawExportSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawExportSettingsImpl struct {
	exportSettings BaseExportSettingsImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawExportSettingsImpl) ExportSettings() BaseExportSettingsImpl {
	return s.exportSettings
}

func (s RawExportSettingsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalExportSettingsImplementation(input []byte) (ExportSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ExportSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLakeExportCommand") {
		var out AzureDatabricksDeltaLakeExportCommand
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeExportCommand: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeExportCopyCommand") {
		var out SnowflakeExportCopyCommand
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeExportCopyCommand: %+v", err)
		}
		return out, nil
	}

	var parent BaseExportSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseExportSettingsImpl: %+v", err)
	}

	return RawExportSettingsImpl{
		exportSettings: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
