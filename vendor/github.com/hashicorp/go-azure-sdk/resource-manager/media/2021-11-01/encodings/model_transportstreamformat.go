package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Format = TransportStreamFormat{}

type TransportStreamFormat struct {
	OutputFiles *[]OutputFile `json:"outputFiles,omitempty"`

	// Fields inherited from Format
	FilenamePattern string `json:"filenamePattern"`
}

var _ json.Marshaler = TransportStreamFormat{}

func (s TransportStreamFormat) MarshalJSON() ([]byte, error) {
	type wrapper TransportStreamFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TransportStreamFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TransportStreamFormat: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.TransportStreamFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TransportStreamFormat: %+v", err)
	}

	return encoded, nil
}
