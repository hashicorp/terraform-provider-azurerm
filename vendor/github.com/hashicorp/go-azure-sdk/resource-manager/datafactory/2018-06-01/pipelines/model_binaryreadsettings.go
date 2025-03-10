package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatReadSettings = BinaryReadSettings{}

type BinaryReadSettings struct {
	CompressionProperties CompressionReadSettings `json:"compressionProperties"`

	// Fields inherited from FormatReadSettings

	Type string `json:"type"`
}

func (s BinaryReadSettings) FormatReadSettings() BaseFormatReadSettingsImpl {
	return BaseFormatReadSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = BinaryReadSettings{}

func (s BinaryReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper BinaryReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BinaryReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BinaryReadSettings: %+v", err)
	}

	decoded["type"] = "BinaryReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BinaryReadSettings: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &BinaryReadSettings{}

func (s *BinaryReadSettings) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BinaryReadSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["compressionProperties"]; ok {
		impl, err := UnmarshalCompressionReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CompressionProperties' for 'BinaryReadSettings': %+v", err)
		}
		s.CompressionProperties = impl
	}

	return nil
}
