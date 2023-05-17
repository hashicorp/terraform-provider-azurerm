package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Preset = StandardEncoderPreset{}

type StandardEncoderPreset struct {
	Codecs              []Codec            `json:"codecs"`
	ExperimentalOptions *map[string]string `json:"experimentalOptions,omitempty"`
	Filters             *Filters           `json:"filters,omitempty"`
	Formats             []Format           `json:"formats"`

	// Fields inherited from Preset
}

var _ json.Marshaler = StandardEncoderPreset{}

func (s StandardEncoderPreset) MarshalJSON() ([]byte, error) {
	type wrapper StandardEncoderPreset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StandardEncoderPreset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StandardEncoderPreset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.StandardEncoderPreset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StandardEncoderPreset: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &StandardEncoderPreset{}

func (s *StandardEncoderPreset) UnmarshalJSON(bytes []byte) error {
	type alias StandardEncoderPreset
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into StandardEncoderPreset: %+v", err)
	}

	s.ExperimentalOptions = decoded.ExperimentalOptions
	s.Filters = decoded.Filters

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StandardEncoderPreset into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["codecs"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Codecs into list []json.RawMessage: %+v", err)
		}

		output := make([]Codec, 0)
		for i, val := range listTemp {
			impl, err := unmarshalCodecImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Codecs' for 'StandardEncoderPreset': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Codecs = output
	}

	if v, ok := temp["formats"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Formats into list []json.RawMessage: %+v", err)
		}

		output := make([]Format, 0)
		for i, val := range listTemp {
			impl, err := unmarshalFormatImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Formats' for 'StandardEncoderPreset': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Formats = output
	}
	return nil
}
