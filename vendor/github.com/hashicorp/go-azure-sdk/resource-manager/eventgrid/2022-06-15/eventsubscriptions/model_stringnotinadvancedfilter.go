package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringNotInAdvancedFilter{}

type StringNotInAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = StringNotInAdvancedFilter{}

func (s StringNotInAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotInAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotInAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotInAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "StringNotIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotInAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
