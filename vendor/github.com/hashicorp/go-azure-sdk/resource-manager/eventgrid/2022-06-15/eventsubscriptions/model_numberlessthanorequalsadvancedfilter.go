package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = NumberLessThanOrEqualsAdvancedFilter{}

type NumberLessThanOrEqualsAdvancedFilter struct {
	Value *float64 `json:"value,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = NumberLessThanOrEqualsAdvancedFilter{}

func (s NumberLessThanOrEqualsAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper NumberLessThanOrEqualsAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NumberLessThanOrEqualsAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NumberLessThanOrEqualsAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "NumberLessThanOrEquals"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NumberLessThanOrEqualsAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
