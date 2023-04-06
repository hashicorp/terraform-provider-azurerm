package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackDescriptor = SelectVideoTrackById{}

type SelectVideoTrackById struct {
	TrackId int64 `json:"trackId"`

	// Fields inherited from TrackDescriptor
}

var _ json.Marshaler = SelectVideoTrackById{}

func (s SelectVideoTrackById) MarshalJSON() ([]byte, error) {
	type wrapper SelectVideoTrackById
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelectVideoTrackById: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelectVideoTrackById: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.SelectVideoTrackById"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelectVideoTrackById: %+v", err)
	}

	return encoded, nil
}
