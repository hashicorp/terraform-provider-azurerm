package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = NumberNotInFilter{}

type NumberNotInFilter struct {
	Values *[]float64 `json:"values,omitempty"`

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s NumberNotInFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = NumberNotInFilter{}

func (s NumberNotInFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberNotInFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberNotInFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberNotInFilter: %+v", err)
	}

	decoded["operatorType"] = "NumberNotIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberNotInFilter: %+v", err)
	}

	return encoded, nil
}
