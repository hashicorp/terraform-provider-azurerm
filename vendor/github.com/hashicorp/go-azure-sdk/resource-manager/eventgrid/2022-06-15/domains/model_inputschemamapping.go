package domains

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputSchemaMapping interface {
}

func unmarshalInputSchemaMappingImplementation(input []byte) (InputSchemaMapping, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InputSchemaMapping into map[string]interface: %+v", err)
	}

	value, ok := temp["inputSchemaMappingType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Json") {
		var out JsonInputSchemaMapping
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonInputSchemaMapping: %+v", err)
		}
		return out, nil
	}

	type RawInputSchemaMappingImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawInputSchemaMappingImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
