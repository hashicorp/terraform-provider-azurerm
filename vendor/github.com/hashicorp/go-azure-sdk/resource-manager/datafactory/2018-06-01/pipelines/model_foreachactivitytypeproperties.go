package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForEachActivityTypeProperties struct {
	Activities   []Activity `json:"activities"`
	BatchCount   *int64     `json:"batchCount,omitempty"`
	IsSequential *bool      `json:"isSequential,omitempty"`
	Items        Expression `json:"items"`
}

var _ json.Unmarshaler = &ForEachActivityTypeProperties{}

func (s *ForEachActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		BatchCount   *int64     `json:"batchCount,omitempty"`
		IsSequential *bool      `json:"isSequential,omitempty"`
		Items        Expression `json:"items"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BatchCount = decoded.BatchCount
	s.IsSequential = decoded.IsSequential
	s.Items = decoded.Items

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ForEachActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["activities"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Activities into list []json.RawMessage: %+v", err)
		}

		output := make([]Activity, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActivityImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Activities' for 'ForEachActivityTypeProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Activities = output
	}

	return nil
}
