package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackDescriptor = SelectVideoTrackByAttribute{}

type SelectVideoTrackByAttribute struct {
	Attribute   TrackAttribute  `json:"attribute"`
	Filter      AttributeFilter `json:"filter"`
	FilterValue *string         `json:"filterValue,omitempty"`

	// Fields inherited from TrackDescriptor
}

var _ json.Marshaler = SelectVideoTrackByAttribute{}

func (s SelectVideoTrackByAttribute) MarshalJSON() ([]byte, error) {
	type wrapper SelectVideoTrackByAttribute
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelectVideoTrackByAttribute: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelectVideoTrackByAttribute: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.SelectVideoTrackByAttribute"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelectVideoTrackByAttribute: %+v", err)
	}

	return encoded, nil
}
