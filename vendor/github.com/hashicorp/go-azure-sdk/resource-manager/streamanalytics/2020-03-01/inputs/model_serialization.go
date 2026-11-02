package inputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Serialization interface {
	Serialization() BaseSerializationImpl
}

var _ Serialization = BaseSerializationImpl{}

type BaseSerializationImpl struct {
	Type EventSerializationType `json:"type"`
}

func (s BaseSerializationImpl) Serialization() BaseSerializationImpl {
	return s
}

var _ Serialization = RawSerializationImpl{}

// RawSerializationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSerializationImpl struct {
	serialization BaseSerializationImpl
	Type          string
	Values        map[string]interface{}
}

func (s RawSerializationImpl) Serialization() BaseSerializationImpl {
	return s.serialization
}

func (s RawSerializationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSerializationImplementation(input []byte) (Serialization, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Serialization into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseSerializationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSerializationImpl: %+v", err)
	}

	return RawSerializationImpl{
		serialization: parent,
		Type:          value,
		Values:        temp,
	}, nil

}
