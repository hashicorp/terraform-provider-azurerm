package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberNotInRangeAdvancedFilter{}

type NumberNotInRangeAdvancedFilter struct {
	Values *[][]float64 `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter

	Key          *string                    `json:"key,omitempty"`
	OperatorType AdvancedFilterOperatorType `json:"operatorType"`
}

func (s NumberNotInRangeAdvancedFilter) AdvancedFilter() BaseAdvancedFilterImpl {
	return BaseAdvancedFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = NumberNotInRangeAdvancedFilter{}

func (s NumberNotInRangeAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberNotInRangeAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberNotInRangeAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberNotInRangeAdvancedFilter: %+v", err)
	}

	decoded["operatorType"] = "NumberNotInRange"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberNotInRangeAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
