package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = H265Video{}

type H265Video struct {
	Complexity           *H265Complexity `json:"complexity,omitempty"`
	KeyFrameInterval     *string         `json:"keyFrameInterval,omitempty"`
	Layers               *[]Layer        `json:"layers,omitempty"`
	SceneChangeDetection *bool           `json:"sceneChangeDetection,omitempty"`
	StretchMode          *StretchMode    `json:"stretchMode,omitempty"`
	SyncMode             *VideoSyncMode  `json:"syncMode,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = H265Video{}

func (s H265Video) MarshalJSON() ([]byte, error) {
	type wrapper H265Video
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling H265Video: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling H265Video: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.H265Video"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling H265Video: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &H265Video{}

func (s *H265Video) UnmarshalJSON(bytes []byte) error {
	type alias H265Video
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into H265Video: %+v", err)
	}

	s.Complexity = decoded.Complexity
	s.KeyFrameInterval = decoded.KeyFrameInterval
	s.Label = decoded.Label
	s.SceneChangeDetection = decoded.SceneChangeDetection
	s.StretchMode = decoded.StretchMode
	s.SyncMode = decoded.SyncMode

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling H265Video into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Layers' for 'H265Video': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Layers = &output
	}
	return nil
}
