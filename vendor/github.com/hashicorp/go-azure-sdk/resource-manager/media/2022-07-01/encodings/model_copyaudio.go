package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Codec = CopyAudio{}

type CopyAudio struct {

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = CopyAudio{}

func (s CopyAudio) MarshalJSON() ([]byte, error) {
	type wrapper CopyAudio
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CopyAudio: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CopyAudio: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.CopyAudio"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CopyAudio: %+v", err)
	}

	return encoded, nil
}
