package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = NumberLessThanFilter{}

type NumberLessThanFilter struct {
	Value *float64 `json:"value,omitempty"`

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s NumberLessThanFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = NumberLessThanFilter{}

func (s NumberLessThanFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberLessThanFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberLessThanFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberLessThanFilter: %+v", err)
	}

	decoded["operatorType"] = "NumberLessThan"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberLessThanFilter: %+v", err)
	}

	return encoded, nil
}
