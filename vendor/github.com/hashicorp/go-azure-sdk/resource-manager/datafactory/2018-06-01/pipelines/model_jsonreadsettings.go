package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatReadSettings = JsonReadSettings{}

type JsonReadSettings struct {
	CompressionProperties CompressionReadSettings `json:"compressionProperties"`

	// Fields inherited from FormatReadSettings

	Type string `json:"type"`
}

func (s JsonReadSettings) FormatReadSettings() BaseFormatReadSettingsImpl {
	return BaseFormatReadSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = JsonReadSettings{}

func (s JsonReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper JsonReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JsonReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JsonReadSettings: %+v", err)
	}

	decoded["type"] = "JsonReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JsonReadSettings: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JsonReadSettings{}

func (s *JsonReadSettings) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JsonReadSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["compressionProperties"]; ok {
		impl, err := UnmarshalCompressionReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CompressionProperties' for 'JsonReadSettings': %+v", err)
		}
		s.CompressionProperties = impl
	}

	return nil
}
