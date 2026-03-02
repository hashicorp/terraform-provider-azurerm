package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UntilActivityTypeProperties struct {
	Activities []Activity   `json:"activities"`
	Expression Expression   `json:"expression"`
	Timeout    *interface{} `json:"timeout,omitempty"`
}

var _ json.Unmarshaler = &UntilActivityTypeProperties{}

func (s *UntilActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Expression Expression   `json:"expression"`
		Timeout    *interface{} `json:"timeout,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Expression = decoded.Expression
	s.Timeout = decoded.Timeout

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UntilActivityTypeProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Activities' for 'UntilActivityTypeProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Activities = output
	}

	return nil
}
