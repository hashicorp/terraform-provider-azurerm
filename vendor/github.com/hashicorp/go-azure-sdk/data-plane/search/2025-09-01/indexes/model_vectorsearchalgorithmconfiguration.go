package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorSearchAlgorithmConfiguration interface {
	VectorSearchAlgorithmConfiguration() BaseVectorSearchAlgorithmConfigurationImpl
}

var _ VectorSearchAlgorithmConfiguration = BaseVectorSearchAlgorithmConfigurationImpl{}

type BaseVectorSearchAlgorithmConfigurationImpl struct {
	Kind VectorSearchAlgorithmKind `json:"kind"`
	Name string                    `json:"name"`
}

func (s BaseVectorSearchAlgorithmConfigurationImpl) VectorSearchAlgorithmConfiguration() BaseVectorSearchAlgorithmConfigurationImpl {
	return s
}

var _ VectorSearchAlgorithmConfiguration = RawVectorSearchAlgorithmConfigurationImpl{}

// RawVectorSearchAlgorithmConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawVectorSearchAlgorithmConfigurationImpl struct {
	vectorSearchAlgorithmConfiguration BaseVectorSearchAlgorithmConfigurationImpl
	Type                               string
	Values                             map[string]interface{}
}

func (s RawVectorSearchAlgorithmConfigurationImpl) VectorSearchAlgorithmConfiguration() BaseVectorSearchAlgorithmConfigurationImpl {
	return s.vectorSearchAlgorithmConfiguration
}

func UnmarshalVectorSearchAlgorithmConfigurationImplementation(input []byte) (VectorSearchAlgorithmConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling VectorSearchAlgorithmConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "exhaustiveKnn") {
		var out ExhaustiveKnnVectorSearchAlgorithmConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExhaustiveKnnVectorSearchAlgorithmConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "hnsw") {
		var out HnswVectorSearchAlgorithmConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HnswVectorSearchAlgorithmConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseVectorSearchAlgorithmConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseVectorSearchAlgorithmConfigurationImpl: %+v", err)
	}

	return RawVectorSearchAlgorithmConfigurationImpl{
		vectorSearchAlgorithmConfiguration: parent,
		Type:                               value,
		Values:                             temp,
	}, nil

}
