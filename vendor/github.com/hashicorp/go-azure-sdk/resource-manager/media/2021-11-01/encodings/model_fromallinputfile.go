package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InputDefinition = FromAllInputFile{}

type FromAllInputFile struct {

	// Fields inherited from InputDefinition
	IncludedTracks *[]TrackDescriptor `json:"includedTracks,omitempty"`
}

var _ json.Marshaler = FromAllInputFile{}

func (s FromAllInputFile) MarshalJSON() ([]byte, error) {
	type wrapper FromAllInputFile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FromAllInputFile: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FromAllInputFile: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.FromAllInputFile"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FromAllInputFile: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &FromAllInputFile{}

func (s *FromAllInputFile) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FromAllInputFile into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["includedTracks"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling IncludedTracks into list []json.RawMessage: %+v", err)
		}

		output := make([]TrackDescriptor, 0)
		for i, val := range listTemp {
			impl, err := unmarshalTrackDescriptorImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'IncludedTracks' for 'FromAllInputFile': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.IncludedTracks = &output
	}
	return nil
}
