package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberGreaterThanOrEqualsAdvancedFilter{}

type NumberGreaterThanOrEqualsAdvancedFilter struct {
	Value *float64 `json:"value,omitempty"`

	// Fields inherited from AdvancedFilter

	Key          *string                    `json:"key,omitempty"`
	OperatorType AdvancedFilterOperatorType `json:"operatorType"`
}

func (s NumberGreaterThanOrEqualsAdvancedFilter) AdvancedFilter() BaseAdvancedFilterImpl {
	return BaseAdvancedFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = NumberGreaterThanOrEqualsAdvancedFilter{}

func (s NumberGreaterThanOrEqualsAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberGreaterThanOrEqualsAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberGreaterThanOrEqualsAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberGreaterThanOrEqualsAdvancedFilter: %+v", err)
	}

	decoded["operatorType"] = "NumberGreaterThanOrEquals"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberGreaterThanOrEqualsAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
