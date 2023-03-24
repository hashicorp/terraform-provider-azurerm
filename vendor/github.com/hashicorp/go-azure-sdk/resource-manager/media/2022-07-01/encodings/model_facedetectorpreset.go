package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Preset = FaceDetectorPreset{}

type FaceDetectorPreset struct {
	BlurType            *BlurType           `json:"blurType,omitempty"`
	ExperimentalOptions *map[string]string  `json:"experimentalOptions,omitempty"`
	Mode                *FaceRedactorMode   `json:"mode,omitempty"`
	Resolution          *AnalysisResolution `json:"resolution,omitempty"`

	// Fields inherited from Preset
}

var _ json.Marshaler = FaceDetectorPreset{}

func (s FaceDetectorPreset) MarshalJSON() ([]byte, error) {
	type wrapper FaceDetectorPreset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FaceDetectorPreset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FaceDetectorPreset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.FaceDetectorPreset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FaceDetectorPreset: %+v", err)
	}

	return encoded, nil
}
