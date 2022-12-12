package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Layer = PngLayer{}

type PngLayer struct {

	// Fields inherited from Layer
	Height *string `json:"height,omitempty"`
	Label  *string `json:"label,omitempty"`
	Width  *string `json:"width,omitempty"`
}

var _ json.Marshaler = PngLayer{}

func (s PngLayer) MarshalJSON() ([]byte, error) {
	type wrapper PngLayer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PngLayer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PngLayer: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.PngLayer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PngLayer: %+v", err)
	}

	return encoded, nil
}
