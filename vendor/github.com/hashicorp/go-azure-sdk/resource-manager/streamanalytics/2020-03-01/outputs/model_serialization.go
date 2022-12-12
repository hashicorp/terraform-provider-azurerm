package outputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Serialization interface {
}

func unmarshalSerializationImplementation(input []byte) (Serialization, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Serialization into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Avro") {
		var out AvroSerialization
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AvroSerialization: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Csv") {
		var out CsvSerialization
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CsvSerialization: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Json") {
		var out JsonSerialization
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonSerialization: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Parquet") {
		var out ParquetSerialization
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetSerialization: %+v", err)
		}
		return out, nil
	}

	type RawSerializationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSerializationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
