package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Preset = BuiltInStandardEncoderPreset{}

type BuiltInStandardEncoderPreset struct {
	Configurations *PresetConfigurations `json:"configurations,omitempty"`
	PresetName     EncoderNamedPreset    `json:"presetName"`

	// Fields inherited from Preset
}

var _ json.Marshaler = BuiltInStandardEncoderPreset{}

func (s BuiltInStandardEncoderPreset) MarshalJSON() ([]byte, error) {
	type wrapper BuiltInStandardEncoderPreset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BuiltInStandardEncoderPreset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BuiltInStandardEncoderPreset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.BuiltInStandardEncoderPreset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BuiltInStandardEncoderPreset: %+v", err)
	}

	return encoded, nil
}
