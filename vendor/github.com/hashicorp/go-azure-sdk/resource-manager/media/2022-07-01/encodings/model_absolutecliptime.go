package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClipTime = AbsoluteClipTime{}

type AbsoluteClipTime struct {
	Time string `json:"time"`

	// Fields inherited from ClipTime
}

var _ json.Marshaler = AbsoluteClipTime{}

func (s AbsoluteClipTime) MarshalJSON() ([]byte, error) {
	type wrapper AbsoluteClipTime
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AbsoluteClipTime: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AbsoluteClipTime: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AbsoluteClipTime"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AbsoluteClipTime: %+v", err)
	}

	return encoded, nil
}
