package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Layer = H265Layer{}

type H265Layer struct {
	AdaptiveBFrame  *bool             `json:"adaptiveBFrame,omitempty"`
	BFrames         *int64            `json:"bFrames,omitempty"`
	Bitrate         int64             `json:"bitrate"`
	BufferWindow    *string           `json:"bufferWindow,omitempty"`
	FrameRate       *string           `json:"frameRate,omitempty"`
	Level           *string           `json:"level,omitempty"`
	MaxBitrate      *int64            `json:"maxBitrate,omitempty"`
	Profile         *H265VideoProfile `json:"profile,omitempty"`
	ReferenceFrames *int64            `json:"referenceFrames,omitempty"`
	Slices          *int64            `json:"slices,omitempty"`

	// Fields inherited from Layer
	Height *string `json:"height,omitempty"`
	Label  *string `json:"label,omitempty"`
	Width  *string `json:"width,omitempty"`
}

var _ json.Marshaler = H265Layer{}

func (s H265Layer) MarshalJSON() ([]byte, error) {
	type wrapper H265Layer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling H265Layer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling H265Layer: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.H265Layer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling H265Layer: %+v", err)
	}

	return encoded, nil
}
