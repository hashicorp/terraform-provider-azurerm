package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = IsNullOrUndefinedFilter{}

type IsNullOrUndefinedFilter struct {

	// Fields inherited from Filter

	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s IsNullOrUndefinedFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Key:          s.Key,
		OperatorType: s.OperatorType,
	}
}

var _ json.Marshaler = IsNullOrUndefinedFilter{}

func (s IsNullOrUndefinedFilter) MarshalJSON() ([]byte, error) {
	type wrapper IsNullOrUndefinedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IsNullOrUndefinedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IsNullOrUndefinedFilter: %+v", err)
	}

	decoded["operatorType"] = "IsNullOrUndefined"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IsNullOrUndefinedFilter: %+v", err)
	}

	return encoded, nil
}
