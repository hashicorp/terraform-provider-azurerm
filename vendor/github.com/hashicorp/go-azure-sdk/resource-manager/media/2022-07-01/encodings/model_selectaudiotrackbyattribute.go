package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackDescriptor = SelectAudioTrackByAttribute{}

type SelectAudioTrackByAttribute struct {
	Attribute      TrackAttribute  `json:"attribute"`
	ChannelMapping *ChannelMapping `json:"channelMapping,omitempty"`
	Filter         AttributeFilter `json:"filter"`
	FilterValue    *string         `json:"filterValue,omitempty"`

	// Fields inherited from TrackDescriptor
}

var _ json.Marshaler = SelectAudioTrackByAttribute{}

func (s SelectAudioTrackByAttribute) MarshalJSON() ([]byte, error) {
	type wrapper SelectAudioTrackByAttribute
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelectAudioTrackByAttribute: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelectAudioTrackByAttribute: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.SelectAudioTrackByAttribute"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelectAudioTrackByAttribute: %+v", err)
	}

	return encoded, nil
}
