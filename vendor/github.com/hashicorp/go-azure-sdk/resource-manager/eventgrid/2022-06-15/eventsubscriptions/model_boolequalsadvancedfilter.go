package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = BoolEqualsAdvancedFilter{}

type BoolEqualsAdvancedFilter struct {
	Value *bool `json:"value,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = BoolEqualsAdvancedFilter{}

func (s BoolEqualsAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper BoolEqualsAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BoolEqualsAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BoolEqualsAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "BoolEquals"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BoolEqualsAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
