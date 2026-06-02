package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FormatReadSettings interface {
	FormatReadSettings() BaseFormatReadSettingsImpl
}

var _ FormatReadSettings = BaseFormatReadSettingsImpl{}

type BaseFormatReadSettingsImpl struct {
	Type string `json:"type"`
}

func (s BaseFormatReadSettingsImpl) FormatReadSettings() BaseFormatReadSettingsImpl {
	return s
}

var _ FormatReadSettings = RawFormatReadSettingsImpl{}

// RawFormatReadSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFormatReadSettingsImpl struct {
	formatReadSettings BaseFormatReadSettingsImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawFormatReadSettingsImpl) FormatReadSettings() BaseFormatReadSettingsImpl {
	return s.formatReadSettings
}

func UnmarshalFormatReadSettingsImplementation(input []byte) (FormatReadSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FormatReadSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BinaryReadSettings") {
		var out BinaryReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BinaryReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DelimitedTextReadSettings") {
		var out DelimitedTextReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DelimitedTextReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JsonReadSettings") {
		var out JsonReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ParquetReadSettings") {
		var out ParquetReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "XmlReadSettings") {
		var out XmlReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XmlReadSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseFormatReadSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFormatReadSettingsImpl: %+v", err)
	}

	return RawFormatReadSettingsImpl{
		formatReadSettings: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
