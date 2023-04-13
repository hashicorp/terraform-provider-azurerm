package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = AacAudio{}

type AacAudio struct {
	Bitrate      *int64           `json:"bitrate,omitempty"`
	Channels     *int64           `json:"channels,omitempty"`
	Profile      *AacAudioProfile `json:"profile,omitempty"`
	SamplingRate *int64           `json:"samplingRate,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = AacAudio{}

func (s AacAudio) MarshalJSON() ([]byte, error) {
	type wrapper AacAudio
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AacAudio: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AacAudio: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AacAudio"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AacAudio: %+v", err)
	}

	return encoded, nil
}
