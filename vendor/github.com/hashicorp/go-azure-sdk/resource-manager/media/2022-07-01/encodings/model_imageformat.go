package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Format = ImageFormat{}

type ImageFormat struct {

	// Fields inherited from Format
	FilenamePattern string `json:"filenamePattern"`
}

var _ json.Marshaler = ImageFormat{}

func (s ImageFormat) MarshalJSON() ([]byte, error) {
	type wrapper ImageFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageFormat: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ImageFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageFormat: %+v", err)
	}

	return encoded, nil
}
