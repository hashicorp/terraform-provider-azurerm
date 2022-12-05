package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Format = MultiBitrateFormat{}

type MultiBitrateFormat struct {
	OutputFiles *[]OutputFile `json:"outputFiles,omitempty"`

	// Fields inherited from Format
	FilenamePattern string `json:"filenamePattern"`
}

var _ json.Marshaler = MultiBitrateFormat{}

func (s MultiBitrateFormat) MarshalJSON() ([]byte, error) {
	type wrapper MultiBitrateFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MultiBitrateFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MultiBitrateFormat: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.MultiBitrateFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MultiBitrateFormat: %+v", err)
	}

	return encoded, nil
}
