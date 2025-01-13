package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberLessThanAdvancedFilter{}

type NumberLessThanAdvancedFilter struct {
	Value *float64 `json:"value,omitempty"`

	// Fields inherited from AdvancedFilter

	Key          *string                    `json:"key,omitempty"`
	OperatorType AdvancedFilterOperatorType `json:"operatorType"`
}

func (s NumberLessThanAdvancedFilter) AdvancedFilter() BaseAdvancedFilterImpl {
	return BaseAdvancedFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = NumberLessThanAdvancedFilter{}

func (s NumberLessThanAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberLessThanAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberLessThanAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberLessThanAdvancedFilter: %+v", err)
	}

	decoded["operatorType"] = "NumberLessThan"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberLessThanAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
