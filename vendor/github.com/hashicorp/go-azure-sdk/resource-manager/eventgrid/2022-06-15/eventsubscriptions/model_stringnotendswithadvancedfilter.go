package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringNotEndsWithAdvancedFilter{}

type StringNotEndsWithAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter

	Key          *string                    `json:"key,omitempty"`
	OperatorType AdvancedFilterOperatorType `json:"operatorType"`
}

func (s StringNotEndsWithAdvancedFilter) AdvancedFilter() BaseAdvancedFilterImpl {
	return BaseAdvancedFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = StringNotEndsWithAdvancedFilter{}

func (s StringNotEndsWithAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotEndsWithAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotEndsWithAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotEndsWithAdvancedFilter: %+v", err)
	}

	decoded["operatorType"] = "StringNotEndsWith"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotEndsWithAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
