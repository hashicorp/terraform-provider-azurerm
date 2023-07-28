package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetLags = CustomTargetLags{}

type CustomTargetLags struct {
	Values []int64 `json:"values"`

	// Fields inherited from TargetLags
}

var _ json.Marshaler = CustomTargetLags{}

func (s CustomTargetLags) MarshalJSON() ([]byte, error) {
	type wrapper CustomTargetLags
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomTargetLags: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomTargetLags: %+v", err)
	}
	decoded["mode"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomTargetLags: %+v", err)
	}

	return encoded, nil
}
