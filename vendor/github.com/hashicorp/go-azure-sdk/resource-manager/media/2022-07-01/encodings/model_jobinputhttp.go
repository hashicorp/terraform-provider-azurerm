package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = JobInputHTTP{}

type JobInputHTTP struct {
	BaseUri          *string            `json:"baseUri,omitempty"`
	End              ClipTime           `json:"end"`
	Files            *[]string          `json:"files,omitempty"`
	InputDefinitions *[]InputDefinition `json:"inputDefinitions,omitempty"`
	Label            *string            `json:"label,omitempty"`
	Start            ClipTime           `json:"start"`

	// Fields inherited from JobInput
}

var _ json.Marshaler = JobInputHTTP{}

func (s JobInputHTTP) MarshalJSON() ([]byte, error) {
	type wrapper JobInputHTTP
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobInputHTTP: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInputHTTP: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobInputHttp"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobInputHTTP: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JobInputHTTP{}

func (s *JobInputHTTP) UnmarshalJSON(bytes []byte) error {
	type alias JobInputHTTP
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JobInputHTTP: %+v", err)
	}

	s.BaseUri = decoded.BaseUri
	s.Files = decoded.Files
	s.Label = decoded.Label

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobInputHTTP into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["end"]; ok {
		impl, err := unmarshalClipTimeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'End' for 'JobInputHTTP': %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'InputDefinitions' for 'JobInputHTTP': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.InputDefinitions = &output
	}

	if v, ok := temp["start"]; ok {
		impl, err := unmarshalClipTimeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Start' for 'JobInputHTTP': %+v", err)
		}
		s.Start = impl
	}
	return nil
}
