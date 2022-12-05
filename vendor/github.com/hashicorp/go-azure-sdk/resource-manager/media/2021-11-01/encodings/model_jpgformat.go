package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Format = JpgFormat{}

type JpgFormat struct {

	// Fields inherited from Format
	FilenamePattern string `json:"filenamePattern"`
}

var _ json.Marshaler = JpgFormat{}

func (s JpgFormat) MarshalJSON() ([]byte, error) {
	type wrapper JpgFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JpgFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JpgFormat: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JpgFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JpgFormat: %+v", err)
	}

	return encoded, nil
}
