package topics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputSchemaMapping interface {
	InputSchemaMapping() BaseInputSchemaMappingImpl
}

var _ InputSchemaMapping = BaseInputSchemaMappingImpl{}

type BaseInputSchemaMappingImpl struct {
	InputSchemaMappingType InputSchemaMappingType `json:"inputSchemaMappingType"`
}

func (s BaseInputSchemaMappingImpl) InputSchemaMapping() BaseInputSchemaMappingImpl {
	return s
}

var _ InputSchemaMapping = RawInputSchemaMappingImpl{}

// RawInputSchemaMappingImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawInputSchemaMappingImpl struct {
	inputSchemaMapping BaseInputSchemaMappingImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawInputSchemaMappingImpl) InputSchemaMapping() BaseInputSchemaMappingImpl {
	return s.inputSchemaMapping
}

func UnmarshalInputSchemaMappingImplementation(input []byte) (InputSchemaMapping, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InputSchemaMapping into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["inputSchemaMappingType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Json") {
		var out JsonInputSchemaMapping
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonInputSchemaMapping: %+v", err)
		}
		return out, nil
	}

	var parent BaseInputSchemaMappingImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseInputSchemaMappingImpl: %+v", err)
	}

	return RawInputSchemaMappingImpl{
		inputSchemaMapping: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
