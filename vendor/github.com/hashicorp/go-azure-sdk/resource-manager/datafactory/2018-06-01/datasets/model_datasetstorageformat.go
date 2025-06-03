package datasets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatasetStorageFormat interface {
	DatasetStorageFormat() BaseDatasetStorageFormatImpl
}

var _ DatasetStorageFormat = BaseDatasetStorageFormatImpl{}

type BaseDatasetStorageFormatImpl struct {
	Deserializer *string `json:"deserializer,omitempty"`
	Serializer   *string `json:"serializer,omitempty"`
	Type         string  `json:"type"`
}

func (s BaseDatasetStorageFormatImpl) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return s
}

var _ DatasetStorageFormat = RawDatasetStorageFormatImpl{}

// RawDatasetStorageFormatImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatasetStorageFormatImpl struct {
	datasetStorageFormat BaseDatasetStorageFormatImpl
	Type                 string
	Values               map[string]interface{}
}

func (s RawDatasetStorageFormatImpl) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return s.datasetStorageFormat
}

func UnmarshalDatasetStorageFormatImplementation(input []byte) (DatasetStorageFormat, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DatasetStorageFormat into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AvroFormat") {
		var out AvroFormat
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AvroFormat: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JsonFormat") {
		var out JsonFormat
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonFormat: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OrcFormat") {
		var out OrcFormat
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OrcFormat: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ParquetFormat") {
		var out ParquetFormat
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetFormat: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TextFormat") {
		var out TextFormat
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TextFormat: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatasetStorageFormatImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatasetStorageFormatImpl: %+v", err)
	}

	return RawDatasetStorageFormatImpl{
		datasetStorageFormat: parent,
		Type:                 value,
		Values:               temp,
	}, nil

}
