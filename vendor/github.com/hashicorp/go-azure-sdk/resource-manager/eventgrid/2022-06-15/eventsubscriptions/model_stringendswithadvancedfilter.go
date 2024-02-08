package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringEndsWithAdvancedFilter{}

type StringEndsWithAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = StringEndsWithAdvancedFilter{}

func (s StringEndsWithAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringEndsWithAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringEndsWithAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringEndsWithAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "StringEndsWith"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringEndsWithAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
