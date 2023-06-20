package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Overlay = AudioOverlay{}

type AudioOverlay struct {

	// Fields inherited from Overlay
	AudioGainLevel  *float64 `json:"audioGainLevel,omitempty"`
	End             *string  `json:"end,omitempty"`
	FadeInDuration  *string  `json:"fadeInDuration,omitempty"`
	FadeOutDuration *string  `json:"fadeOutDuration,omitempty"`
	InputLabel      string   `json:"inputLabel"`
	Start           *string  `json:"start,omitempty"`
}

var _ json.Marshaler = AudioOverlay{}

func (s AudioOverlay) MarshalJSON() ([]byte, error) {
	type wrapper AudioOverlay
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AudioOverlay: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AudioOverlay: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AudioOverlay"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AudioOverlay: %+v", err)
	}

	return encoded, nil
}
