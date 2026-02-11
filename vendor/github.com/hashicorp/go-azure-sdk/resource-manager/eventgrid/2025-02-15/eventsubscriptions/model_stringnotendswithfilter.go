package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = StringNotEndsWithFilter{}

type StringNotEndsWithFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s StringNotEndsWithFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = StringNotEndsWithFilter{}

func (s StringNotEndsWithFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotEndsWithFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotEndsWithFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotEndsWithFilter: %+v", err)
	}

	decoded["operatorType"] = "StringNotEndsWith"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotEndsWithFilter: %+v", err)
	}

	return encoded, nil
}
