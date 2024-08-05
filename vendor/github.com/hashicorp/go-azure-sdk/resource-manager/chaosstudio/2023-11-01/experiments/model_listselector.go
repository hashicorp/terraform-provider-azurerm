package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Selector = ListSelector{}

type ListSelector struct {
	Targets []TargetReference `json:"targets"`

	// Fields inherited from Selector
	Filter Filter `json:"filter"`
	Id     string `json:"id"`
}

var _ json.Marshaler = ListSelector{}

func (s ListSelector) MarshalJSON() ([]byte, error) {
	type wrapper ListSelector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ListSelector: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ListSelector: %+v", err)
	}
	decoded["type"] = "List"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ListSelector: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ListSelector{}

func (s *ListSelector) UnmarshalJSON(bytes []byte) error {
	type alias ListSelector
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ListSelector: %+v", err)
	}

	s.Id = decoded.Id
	s.Targets = decoded.Targets

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ListSelector into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["filter"]; ok {
		impl, err := unmarshalFilterImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Filter' for 'ListSelector': %+v", err)
		}
		s.Filter = impl
	}
	return nil
}
