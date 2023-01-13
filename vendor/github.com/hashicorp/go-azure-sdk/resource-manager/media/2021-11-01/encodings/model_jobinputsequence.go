package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = JobInputSequence{}

type JobInputSequence struct {
	Inputs *[]JobInput `json:"inputs,omitempty"`

	// Fields inherited from JobInput
}

var _ json.Marshaler = JobInputSequence{}

func (s JobInputSequence) MarshalJSON() ([]byte, error) {
	type wrapper JobInputSequence
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobInputSequence: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInputSequence: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobInputSequence"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobInputSequence: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JobInputSequence{}

func (s *JobInputSequence) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobInputSequence into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Inputs' for 'JobInputSequence': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Inputs = &output
	}
	return nil
}
