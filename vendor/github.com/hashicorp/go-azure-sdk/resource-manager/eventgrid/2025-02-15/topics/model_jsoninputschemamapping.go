package topics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InputSchemaMapping = JsonInputSchemaMapping{}

type JsonInputSchemaMapping struct {
	Properties *JsonInputSchemaMappingProperties `json:"properties,omitempty"`

	// Fields inherited from InputSchemaMapping

	InputSchemaMappingType InputSchemaMappingType `json:"inputSchemaMappingType"`
}

func (s JsonInputSchemaMapping) InputSchemaMapping() BaseInputSchemaMappingImpl {
	return BaseInputSchemaMappingImpl{
		InputSchemaMappingType: s.InputSchemaMappingType,
	}
}

var _ json.Marshaler = JsonInputSchemaMapping{}

func (s JsonInputSchemaMapping) MarshalJSON() ([]byte, error) {
	type wrapper JsonInputSchemaMapping
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JsonInputSchemaMapping: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JsonInputSchemaMapping: %+v", err)
	}

	decoded["inputSchemaMappingType"] = "Json"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JsonInputSchemaMapping: %+v", err)
	}

	return encoded, nil
}
