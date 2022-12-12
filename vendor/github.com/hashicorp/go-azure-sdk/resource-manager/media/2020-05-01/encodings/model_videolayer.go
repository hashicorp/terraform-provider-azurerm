package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Layer = VideoLayer{}

type VideoLayer struct {
	AdaptiveBFrame *bool   `json:"adaptiveBFrame,omitempty"`
	BFrames        *int64  `json:"bFrames,omitempty"`
	Bitrate        int64   `json:"bitrate"`
	FrameRate      *string `json:"frameRate,omitempty"`
	MaxBitrate     *int64  `json:"maxBitrate,omitempty"`
	Slices         *int64  `json:"slices,omitempty"`

	// Fields inherited from Layer
	Height *string `json:"height,omitempty"`
	Label  *string `json:"label,omitempty"`
	Width  *string `json:"width,omitempty"`
}

var _ json.Marshaler = VideoLayer{}

func (s VideoLayer) MarshalJSON() ([]byte, error) {
	type wrapper VideoLayer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VideoLayer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VideoLayer: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.VideoLayer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VideoLayer: %+v", err)
	}

	return encoded, nil
}
