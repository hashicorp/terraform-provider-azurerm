package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Overlay = VideoOverlay{}

type VideoOverlay struct {
	CropRectangle *Rectangle `json:"cropRectangle,omitempty"`
	Opacity       *float64   `json:"opacity,omitempty"`
	Position      *Rectangle `json:"position,omitempty"`

	// Fields inherited from Overlay
	AudioGainLevel  *float64 `json:"audioGainLevel,omitempty"`
	End             *string  `json:"end,omitempty"`
	FadeInDuration  *string  `json:"fadeInDuration,omitempty"`
	FadeOutDuration *string  `json:"fadeOutDuration,omitempty"`
	InputLabel      string   `json:"inputLabel"`
	Start           *string  `json:"start,omitempty"`
}

var _ json.Marshaler = VideoOverlay{}

func (s VideoOverlay) MarshalJSON() ([]byte, error) {
	type wrapper VideoOverlay
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VideoOverlay: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VideoOverlay: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.VideoOverlay"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VideoOverlay: %+v", err)
	}

	return encoded, nil
}
