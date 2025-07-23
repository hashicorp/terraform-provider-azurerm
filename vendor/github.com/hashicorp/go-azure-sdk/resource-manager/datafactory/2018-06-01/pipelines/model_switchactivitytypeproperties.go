package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchActivityTypeProperties struct {
	Cases             *[]SwitchCase `json:"cases,omitempty"`
	DefaultActivities *[]Activity   `json:"defaultActivities,omitempty"`
	On                Expression    `json:"on"`
}

var _ json.Unmarshaler = &SwitchActivityTypeProperties{}

func (s *SwitchActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Cases *[]SwitchCase `json:"cases,omitempty"`
		On    Expression    `json:"on"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Cases = decoded.Cases
	s.On = decoded.On

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SwitchActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["defaultActivities"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling DefaultActivities into list []json.RawMessage: %+v", err)
		}

		output := make([]Activity, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActivityImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'DefaultActivities' for 'SwitchActivityTypeProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.DefaultActivities = &output
	}

	return nil
}
