package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Layer = H265VideoLayer{}

type H265VideoLayer struct {
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

var _ json.Marshaler = H265VideoLayer{}

func (s H265VideoLayer) MarshalJSON() ([]byte, error) {
	type wrapper H265VideoLayer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling H265VideoLayer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling H265VideoLayer: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.H265VideoLayer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling H265VideoLayer: %+v", err)
	}

	return encoded, nil
}
