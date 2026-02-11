package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = StringNotBeginsWithFilter{}

type StringNotBeginsWithFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s StringNotBeginsWithFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = StringNotBeginsWithFilter{}

func (s StringNotBeginsWithFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotBeginsWithFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotBeginsWithFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotBeginsWithFilter: %+v", err)
	}

	decoded["operatorType"] = "StringNotBeginsWith"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotBeginsWithFilter: %+v", err)
	}

	return encoded, nil
}
