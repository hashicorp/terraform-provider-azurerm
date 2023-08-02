package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringInAdvancedFilter{}

type StringInAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = StringInAdvancedFilter{}

func (s StringInAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringInAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringInAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringInAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "StringIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringInAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
