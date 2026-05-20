package streamingjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputProperties interface {
	InputProperties() BaseInputPropertiesImpl
}

var _ InputProperties = BaseInputPropertiesImpl{}

type BaseInputPropertiesImpl struct {
	Compression   *Compression  `json:"compression,omitempty"`
	Diagnostics   *Diagnostics  `json:"diagnostics,omitempty"`
	Etag          *string       `json:"etag,omitempty"`
	PartitionKey  *string       `json:"partitionKey,omitempty"`
	Serialization Serialization `json:"serialization"`
	Type          string        `json:"type"`
}

func (s BaseInputPropertiesImpl) InputProperties() BaseInputPropertiesImpl {
	return s
}

var _ InputProperties = RawInputPropertiesImpl{}

// RawInputPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawInputPropertiesImpl struct {
	inputProperties BaseInputPropertiesImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawInputPropertiesImpl) InputProperties() BaseInputPropertiesImpl {
	return s.inputProperties
}

var _ json.Unmarshaler = &BaseInputPropertiesImpl{}

func (s *BaseInputPropertiesImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Compression  *Compression `json:"compression,omitempty"`
		Diagnostics  *Diagnostics `json:"diagnostics,omitempty"`
		Etag         *string      `json:"etag,omitempty"`
		PartitionKey *string      `json:"partitionKey,omitempty"`
		Type         string       `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Compression = decoded.Compression
	s.Diagnostics = decoded.Diagnostics
	s.Etag = decoded.Etag
	s.PartitionKey = decoded.PartitionKey
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseInputPropertiesImpl into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["serialization"]; ok {
		impl, err := UnmarshalSerializationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Serialization' for 'BaseInputPropertiesImpl': %+v", err)
		}
		s.Serialization = impl
	}

	return nil
}

func UnmarshalInputPropertiesImplementation(input []byte) (InputProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InputProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Reference") {
		var out ReferenceInputProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ReferenceInputProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Stream") {
		var out StreamInputProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StreamInputProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseInputPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseInputPropertiesImpl: %+v", err)
	}

	return RawInputPropertiesImpl{
		inputProperties: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
