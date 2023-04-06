package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackDescriptor = SelectAudioTrackById{}

type SelectAudioTrackById struct {
	ChannelMapping *ChannelMapping `json:"channelMapping,omitempty"`
	TrackId        int64           `json:"trackId"`

	// Fields inherited from TrackDescriptor
}

var _ json.Marshaler = SelectAudioTrackById{}

func (s SelectAudioTrackById) MarshalJSON() ([]byte, error) {
	type wrapper SelectAudioTrackById
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelectAudioTrackById: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelectAudioTrackById: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.SelectAudioTrackById"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelectAudioTrackById: %+v", err)
	}

	return encoded, nil
}
