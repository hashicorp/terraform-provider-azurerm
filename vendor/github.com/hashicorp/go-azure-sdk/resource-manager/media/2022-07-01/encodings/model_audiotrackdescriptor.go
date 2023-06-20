package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackDescriptor = AudioTrackDescriptor{}

type AudioTrackDescriptor struct {
	ChannelMapping *ChannelMapping `json:"channelMapping,omitempty"`

	// Fields inherited from TrackDescriptor
}

var _ json.Marshaler = AudioTrackDescriptor{}

func (s AudioTrackDescriptor) MarshalJSON() ([]byte, error) {
	type wrapper AudioTrackDescriptor
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AudioTrackDescriptor: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AudioTrackDescriptor: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AudioTrackDescriptor"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AudioTrackDescriptor: %+v", err)
	}

	return encoded, nil
}
