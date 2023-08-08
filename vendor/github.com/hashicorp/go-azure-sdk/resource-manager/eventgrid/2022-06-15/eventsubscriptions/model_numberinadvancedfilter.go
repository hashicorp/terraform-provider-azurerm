package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberInAdvancedFilter{}

type NumberInAdvancedFilter struct {
	Values *[]float64 `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = NumberInAdvancedFilter{}

func (s NumberInAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberInAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberInAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberInAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "NumberIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberInAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
