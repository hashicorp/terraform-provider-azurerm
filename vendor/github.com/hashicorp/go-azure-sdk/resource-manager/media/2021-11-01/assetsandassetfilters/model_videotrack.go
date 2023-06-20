package assetsandassetfilters

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackBase = VideoTrack{}

type VideoTrack struct {

	// Fields inherited from TrackBase
}

var _ json.Marshaler = VideoTrack{}

func (s VideoTrack) MarshalJSON() ([]byte, error) {
	type wrapper VideoTrack
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VideoTrack: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VideoTrack: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.VideoTrack"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VideoTrack: %+v", err)
	}

	return encoded, nil
}
