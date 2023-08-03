package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringNotContainsAdvancedFilter{}

type StringNotContainsAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = StringNotContainsAdvancedFilter{}

func (s StringNotContainsAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotContainsAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotContainsAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotContainsAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "StringNotContains"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotContainsAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
