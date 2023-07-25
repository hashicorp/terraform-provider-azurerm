package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = H264Video{}

type H264Video struct {
	Complexity           *H264Complexity      `json:"complexity,omitempty"`
	KeyFrameInterval     *string              `json:"keyFrameInterval,omitempty"`
	Layers               *[]H264Layer         `json:"layers,omitempty"`
	RateControlMode      *H264RateControlMode `json:"rateControlMode,omitempty"`
	SceneChangeDetection *bool                `json:"sceneChangeDetection,omitempty"`
	StretchMode          *StretchMode         `json:"stretchMode,omitempty"`
	SyncMode             *VideoSyncMode       `json:"syncMode,omitempty"`

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
