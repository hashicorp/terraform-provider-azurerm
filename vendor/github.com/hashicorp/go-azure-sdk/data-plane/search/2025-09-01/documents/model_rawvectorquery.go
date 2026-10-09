package documents

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorQuery = RawVectorQuery{}

type RawVectorQuery struct {
	Vector []float64 `json:"vector"`

	// Fields inherited from VectorQuery

	Exhaustive   *bool           `json:"exhaustive,omitempty"`
	Fields       *string         `json:"fields,omitempty"`
	K            *int64          `json:"k,omitempty"`
	Kind         VectorQueryKind `json:"kind"`
	Oversampling *float64        `json:"oversampling,omitempty"`
	Weight       *float64        `json:"weight,omitempty"`
}

func (s RawVectorQuery) VectorQuery() BaseVectorQueryImpl {
	return BaseVectorQueryImpl{
		Exhaustive:   s.Exhaustive,
		Fields:       s.Fields,
		K:            s.K,
		Kind:         s.Kind,
		Oversampling: s.Oversampling,
		Weight:       s.Weight,
	}
}

var _ json.Marshaler = RawVectorQuery{}

func (s RawVectorQuery) MarshalJSON() ([]byte, error) {
	type wrapper RawVectorQuery
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RawVectorQuery: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RawVectorQuery: %+v", err)
	}

	decoded["kind"] = "vector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RawVectorQuery: %+v", err)
	}

	return encoded, nil
}
