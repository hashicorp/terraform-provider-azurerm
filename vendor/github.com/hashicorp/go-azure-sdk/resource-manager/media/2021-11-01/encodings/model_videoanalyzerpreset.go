package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Preset = VideoAnalyzerPreset{}

type VideoAnalyzerPreset struct {
	AudioLanguage       *string            `json:"audioLanguage,omitempty"`
	ExperimentalOptions *map[string]string `json:"experimentalOptions,omitempty"`
	InsightsToExtract   *InsightsType      `json:"insightsToExtract,omitempty"`
	Mode                *AudioAnalysisMode `json:"mode,omitempty"`

	// Fields inherited from Preset
}

var _ json.Marshaler = VideoAnalyzerPreset{}

func (s VideoAnalyzerPreset) MarshalJSON() ([]byte, error) {
	type wrapper VideoAnalyzerPreset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VideoAnalyzerPreset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VideoAnalyzerPreset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.VideoAnalyzerPreset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VideoAnalyzerPreset: %+v", err)
	}

	return encoded, nil
}
