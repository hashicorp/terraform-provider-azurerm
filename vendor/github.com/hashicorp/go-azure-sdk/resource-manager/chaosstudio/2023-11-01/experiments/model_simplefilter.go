package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Filter = SimpleFilter{}

type SimpleFilter struct {
	Parameters *SimpleFilterParameters `json:"parameters,omitempty"`

	// Fields inherited from Filter

	Type FilterType `json:"type"`
}

func (s SimpleFilter) Filter() BaseFilterImpl {
	return BaseFilterImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SimpleFilter{}

func (s SimpleFilter) MarshalJSON() ([]byte, error) {
	type wrapper SimpleFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SimpleFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SimpleFilter: %+v", err)
	}

	decoded["type"] = "Simple"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SimpleFilter: %+v", err)
	}

	return encoded, nil
}
