package documents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorQuery interface {
	VectorQuery() BaseVectorQueryImpl
}

var _ VectorQuery = BaseVectorQueryImpl{}

type BaseVectorQueryImpl struct {
	Exhaustive   *bool           `json:"exhaustive,omitempty"`
	Fields       *string         `json:"fields,omitempty"`
	K            *int64          `json:"k,omitempty"`
	Kind         VectorQueryKind `json:"kind"`
	Oversampling *float64        `json:"oversampling,omitempty"`
	Weight       *float64        `json:"weight,omitempty"`
}

func (s BaseVectorQueryImpl) VectorQuery() BaseVectorQueryImpl {
	return s
}

var _ VectorQuery = RawVectorQueryImpl{}

// RawVectorQueryImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawVectorQueryImpl struct {
	vectorQuery BaseVectorQueryImpl
	Type        string
	Values      map[string]interface{}
}

func (s RawVectorQueryImpl) VectorQuery() BaseVectorQueryImpl {
	return s.vectorQuery
}

func UnmarshalVectorQueryImplementation(input []byte) (VectorQuery, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling VectorQuery into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "vector") {
		var out RawVectorQuery
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RawVectorQuery: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "text") {
		var out VectorizableTextQuery
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VectorizableTextQuery: %+v", err)
		}
		return out, nil
	}

	var parent BaseVectorQueryImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseVectorQueryImpl: %+v", err)
	}

	return RawVectorQueryImpl{
		vectorQuery: parent,
		Type:        value,
		Values:      temp,
	}, nil

}
