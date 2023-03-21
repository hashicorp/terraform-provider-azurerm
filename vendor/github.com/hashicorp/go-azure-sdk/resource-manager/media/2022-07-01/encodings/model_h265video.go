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
	Layers               *[]H265Layer    `json:"layers,omitempty"`
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
