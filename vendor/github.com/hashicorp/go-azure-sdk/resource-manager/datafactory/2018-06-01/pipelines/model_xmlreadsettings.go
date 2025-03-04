package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatReadSettings = XmlReadSettings{}

type XmlReadSettings struct {
	CompressionProperties CompressionReadSettings `json:"compressionProperties"`
	DetectDataType        *bool                   `json:"detectDataType,omitempty"`
	NamespacePrefixes     *map[string]string      `json:"namespacePrefixes,omitempty"`
	Namespaces            *bool                   `json:"namespaces,omitempty"`
	ValidationMode        *interface{}            `json:"validationMode,omitempty"`

	// Fields inherited from FormatReadSettings

	Type string `json:"type"`
}

func (s XmlReadSettings) FormatReadSettings() BaseFormatReadSettingsImpl {
	return BaseFormatReadSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = XmlReadSettings{}

func (s XmlReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper XmlReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling XmlReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling XmlReadSettings: %+v", err)
	}

	decoded["type"] = "XmlReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling XmlReadSettings: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &XmlReadSettings{}

func (s *XmlReadSettings) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DetectDataType    *bool              `json:"detectDataType,omitempty"`
		NamespacePrefixes *map[string]string `json:"namespacePrefixes,omitempty"`
		Namespaces        *bool              `json:"namespaces,omitempty"`
		ValidationMode    *interface{}       `json:"validationMode,omitempty"`
		Type              string             `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DetectDataType = decoded.DetectDataType
	s.NamespacePrefixes = decoded.NamespacePrefixes
	s.Namespaces = decoded.Namespaces
	s.ValidationMode = decoded.ValidationMode
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling XmlReadSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["compressionProperties"]; ok {
		impl, err := UnmarshalCompressionReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CompressionProperties' for 'XmlReadSettings': %+v", err)
		}
		s.CompressionProperties = impl
	}

	return nil
}
