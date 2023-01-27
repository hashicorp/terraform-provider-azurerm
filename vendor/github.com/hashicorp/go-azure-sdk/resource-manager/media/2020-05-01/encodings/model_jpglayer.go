package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Layer = JpgLayer{}

type JpgLayer struct {
	Quality *int64 `json:"quality,omitempty"`

	// Fields inherited from Layer
	Height *string `json:"height,omitempty"`
	Label  *string `json:"label,omitempty"`
	Width  *string `json:"width,omitempty"`
}

var _ json.Marshaler = JpgLayer{}

func (s JpgLayer) MarshalJSON() ([]byte, error) {
	type wrapper JpgLayer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JpgLayer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JpgLayer: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JpgLayer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JpgLayer: %+v", err)
	}

	return encoded, nil
}
