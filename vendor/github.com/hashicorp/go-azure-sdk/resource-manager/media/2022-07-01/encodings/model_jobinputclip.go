package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = JobInputClip{}

type JobInputClip struct {
	End              ClipTime           `json:"end"`
	Files            *[]string          `json:"files,omitempty"`
	InputDefinitions *[]InputDefinition `json:"inputDefinitions,omitempty"`
	Label            *string            `json:"label,omitempty"`
	Start            ClipTime           `json:"start"`

	// Fields inherited from JobInput
}

var _ json.Marshaler = JobInputClip{}

func (s JobInputClip) MarshalJSON() ([]byte, error) {
	type wrapper JobInputClip
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobInputClip: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInputClip: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobInputClip"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobInputClip: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JobInputClip{}

func (s *JobInputClip) UnmarshalJSON(bytes []byte) error {
	type alias JobInputClip
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JobInputClip: %+v", err)
	}

	s.Files = decoded.Files
	s.Label = decoded.Label

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobInputClip into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["end"]; ok {
		impl, err := unmarshalClipTimeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'End' for 'JobInputClip': %+v", err)
		}
		s.End = impl
	}

	if v, ok := temp["inputDefinitions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling InputDefinitions into list []json.RawMessage: %+v", err)
		}

		output := make([]InputDefinition, 0)
		for i, val := range listTemp {
			impl, err := unmarshalInputDefinitionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'InputDefinitions' for 'JobInputClip': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.InputDefinitions = &output
	}

	if v, ok := temp["start"]; ok {
		impl, err := unmarshalClipTimeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Start' for 'JobInputClip': %+v", err)
		}
		s.Start = impl
	}
	return nil
}
