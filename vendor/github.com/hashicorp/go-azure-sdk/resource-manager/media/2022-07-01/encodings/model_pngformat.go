package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Format = PngFormat{}

type PngFormat struct {

	// Fields inherited from Format
	FilenamePattern string `json:"filenamePattern"`
}

var _ json.Marshaler = PngFormat{}

func (s PngFormat) MarshalJSON() ([]byte, error) {
	type wrapper PngFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PngFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PngFormat: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.PngFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PngFormat: %+v", err)
	}

	return encoded, nil
}
