package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = StringInFilter{}

type StringInFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s StringInFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = StringInFilter{}

func (s StringInFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringInFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringInFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringInFilter: %+v", err)
	}

	decoded["operatorType"] = "StringIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringInFilter: %+v", err)
	}

	return encoded, nil
}
