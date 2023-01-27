package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = H264Video{}

type H264Video struct {
	Complexity           *H264Complexity `json:"complexity,omitempty"`
	KeyFrameInterval     *string         `json:"keyFrameInterval,omitempty"`
	Layers               *[]Layer        `json:"layers,omitempty"`
	SceneChangeDetection *bool           `json:"sceneChangeDetection,omitempty"`
	StretchMode          *StretchMode    `json:"stretchMode,omitempty"`
	SyncMode             *VideoSyncMode  `json:"syncMode,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = H264Video{}

func (s H264Video) MarshalJSON() ([]byte, error) {
	type wrapper H264Video
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling H264Video: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling H264Video: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.H264Video"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling H264Video: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &H264Video{}

func (s *H264Video) UnmarshalJSON(bytes []byte) error {
	type alias H264Video
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into H264Video: %+v", err)
	}

	s.Complexity = decoded.Complexity
	s.KeyFrameInterval = decoded.KeyFrameInterval
	s.Label = decoded.Label
	s.SceneChangeDetection = decoded.SceneChangeDetection
	s.StretchMode = decoded.StretchMode
	s.SyncMode = decoded.SyncMode

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling H264Video into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["layers"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Layers into list []json.RawMessage: %+v", err)
		}

		output := make([]Layer, 0)
		for i, val := range listTemp {
			impl, err := unmarshalLayerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Layers' for 'H264Video': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Layers = &output
	}
	return nil
}
