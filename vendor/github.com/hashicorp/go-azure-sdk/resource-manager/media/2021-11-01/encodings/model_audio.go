package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = Audio{}

type Audio struct {
	Bitrate      *int64 `json:"bitrate,omitempty"`
	Channels     *int64 `json:"channels,omitempty"`
	SamplingRate *int64 `json:"samplingRate,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = Audio{}

func (s Audio) MarshalJSON() ([]byte, error) {
	type wrapper Audio
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Audio: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Audio: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.Audio"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Audio: %+v", err)
	}

	return encoded, nil
}
