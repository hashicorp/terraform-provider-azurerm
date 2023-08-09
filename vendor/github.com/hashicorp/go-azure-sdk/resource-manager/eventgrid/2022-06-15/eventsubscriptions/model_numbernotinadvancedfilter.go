package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberNotInAdvancedFilter{}

type NumberNotInAdvancedFilter struct {
	Values *[]float64 `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = NumberNotInAdvancedFilter{}

func (s NumberNotInAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberNotInAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberNotInAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberNotInAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "NumberNotIn"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberNotInAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
