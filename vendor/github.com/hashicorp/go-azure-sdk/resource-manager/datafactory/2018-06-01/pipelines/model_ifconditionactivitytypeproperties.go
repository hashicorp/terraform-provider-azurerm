package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IfConditionActivityTypeProperties struct {
	Expression        Expression  `json:"expression"`
	IfFalseActivities *[]Activity `json:"ifFalseActivities,omitempty"`
	IfTrueActivities  *[]Activity `json:"ifTrueActivities,omitempty"`
}

var _ json.Unmarshaler = &IfConditionActivityTypeProperties{}

func (s *IfConditionActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Expression Expression `json:"expression"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Expression = decoded.Expression

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling IfConditionActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["ifFalseActivities"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling IfFalseActivities into list []json.RawMessage: %+v", err)
		}

		output := make([]Activity, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActivityImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'IfFalseActivities' for 'IfConditionActivityTypeProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.IfFalseActivities = &output
	}

	if v, ok := temp["ifTrueActivities"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling IfTrueActivities into list []json.RawMessage: %+v", err)
		}

		output := make([]Activity, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActivityImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'IfTrueActivities' for 'IfConditionActivityTypeProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.IfTrueActivities = &output
	}

	return nil
}
