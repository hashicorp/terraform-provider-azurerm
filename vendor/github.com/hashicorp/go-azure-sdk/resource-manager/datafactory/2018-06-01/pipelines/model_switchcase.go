package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchCase struct {
	Activities *[]Activity `json:"activities,omitempty"`
	Value      *string     `json:"value,omitempty"`
}

var _ json.Unmarshaler = &SwitchCase{}

func (s *SwitchCase) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Value *string `json:"value,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Value = decoded.Value

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SwitchCase into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Activities' for 'SwitchCase': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Activities = &output
	}

	return nil
}
