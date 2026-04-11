package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorSearchCompressionConfiguration interface {
	VectorSearchCompressionConfiguration() BaseVectorSearchCompressionConfigurationImpl
}

var _ VectorSearchCompressionConfiguration = BaseVectorSearchCompressionConfigurationImpl{}

type BaseVectorSearchCompressionConfigurationImpl struct {
	Kind                VectorSearchCompressionKind `json:"kind"`
	Name                string                      `json:"name"`
	RescoringOptions    *RescoringOptions           `json:"rescoringOptions,omitempty"`
	TruncationDimension *int64                      `json:"truncationDimension,omitempty"`
}

func (s BaseVectorSearchCompressionConfigurationImpl) VectorSearchCompressionConfiguration() BaseVectorSearchCompressionConfigurationImpl {
	return s
}

var _ VectorSearchCompressionConfiguration = RawVectorSearchCompressionConfigurationImpl{}

// RawVectorSearchCompressionConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawVectorSearchCompressionConfigurationImpl struct {
	vectorSearchCompressionConfiguration BaseVectorSearchCompressionConfigurationImpl
	Type                                 string
	Values                               map[string]interface{}
}

func (s RawVectorSearchCompressionConfigurationImpl) VectorSearchCompressionConfiguration() BaseVectorSearchCompressionConfigurationImpl {
	return s.vectorSearchCompressionConfiguration
}

func UnmarshalVectorSearchCompressionConfigurationImplementation(input []byte) (VectorSearchCompressionConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling VectorSearchCompressionConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "binaryQuantization") {
		var out BinaryQuantizationVectorSearchCompressionConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BinaryQuantizationVectorSearchCompressionConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "scalarQuantization") {
		var out ScalarQuantizationVectorSearchCompressionConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScalarQuantizationVectorSearchCompressionConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseVectorSearchCompressionConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseVectorSearchCompressionConfigurationImpl: %+v", err)
	}

	return RawVectorSearchCompressionConfigurationImpl{
		vectorSearchCompressionConfiguration: parent,
		Type:                                 value,
		Values:                               temp,
	}, nil

}
