package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Preset = AudioAnalyzerPreset{}

type AudioAnalyzerPreset struct {
	AudioLanguage       *string            `json:"audioLanguage,omitempty"`
	ExperimentalOptions *map[string]string `json:"experimentalOptions,omitempty"`
	Mode                *AudioAnalysisMode `json:"mode,omitempty"`

	// Fields inherited from Preset
}

var _ json.Marshaler = AudioAnalyzerPreset{}

func (s AudioAnalyzerPreset) MarshalJSON() ([]byte, error) {
	type wrapper AudioAnalyzerPreset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AudioAnalyzerPreset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AudioAnalyzerPreset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AudioAnalyzerPreset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AudioAnalyzerPreset: %+v", err)
	}

	return encoded, nil
}
