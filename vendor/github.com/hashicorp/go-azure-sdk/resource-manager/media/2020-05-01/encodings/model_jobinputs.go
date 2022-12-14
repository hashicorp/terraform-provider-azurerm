package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = JobInputs{}

type JobInputs struct {
	Inputs *[]JobInput `json:"inputs,omitempty"`

	// Fields inherited from JobInput
}

var _ json.Marshaler = JobInputs{}

func (s JobInputs) MarshalJSON() ([]byte, error) {
	type wrapper JobInputs
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobInputs: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInputs: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobInputs"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobInputs: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JobInputs{}

func (s *JobInputs) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobInputs into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["inputs"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Inputs into list []json.RawMessage: %+v", err)
		}

		output := make([]JobInput, 0)
		for i, val := range listTemp {
			impl, err := unmarshalJobInputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Inputs' for 'JobInputs': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Inputs = &output
	}
	return nil
}
