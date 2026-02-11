package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FiltersConfiguration struct {
	Filters            *[]Filter `json:"filters,omitempty"`
	IncludedEventTypes *[]string `json:"includedEventTypes,omitempty"`
}

var _ json.Unmarshaler = &FiltersConfiguration{}

func (s *FiltersConfiguration) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		IncludedEventTypes *[]string `json:"includedEventTypes,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.IncludedEventTypes = decoded.IncludedEventTypes

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FiltersConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["filters"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Filters into list []json.RawMessage: %+v", err)
		}

		output := make([]Filter, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalFilterImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Filters' for 'FiltersConfiguration': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Filters = &output
	}

	return nil
}
