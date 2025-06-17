package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CompressionReadSettings interface {
	CompressionReadSettings() BaseCompressionReadSettingsImpl
}

var _ CompressionReadSettings = BaseCompressionReadSettingsImpl{}

type BaseCompressionReadSettingsImpl struct {
	Type string `json:"type"`
}

func (s BaseCompressionReadSettingsImpl) CompressionReadSettings() BaseCompressionReadSettingsImpl {
	return s
}

var _ CompressionReadSettings = RawCompressionReadSettingsImpl{}

// RawCompressionReadSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCompressionReadSettingsImpl struct {
	compressionReadSettings BaseCompressionReadSettingsImpl
	Type                    string
	Values                  map[string]interface{}
}

func (s RawCompressionReadSettingsImpl) CompressionReadSettings() BaseCompressionReadSettingsImpl {
	return s.compressionReadSettings
}

func UnmarshalCompressionReadSettingsImplementation(input []byte) (CompressionReadSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CompressionReadSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "TarGZipReadSettings") {
		var out TarGZipReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TarGZipReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TarReadSettings") {
		var out TarReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TarReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ZipDeflateReadSettings") {
		var out ZipDeflateReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ZipDeflateReadSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseCompressionReadSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCompressionReadSettingsImpl: %+v", err)
	}

	return RawCompressionReadSettingsImpl{
		compressionReadSettings: parent,
		Type:                    value,
		Values:                  temp,
	}, nil

}
