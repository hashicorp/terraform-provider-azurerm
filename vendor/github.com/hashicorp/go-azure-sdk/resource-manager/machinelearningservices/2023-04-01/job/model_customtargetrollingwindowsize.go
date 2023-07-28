package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetRollingWindowSize = CustomTargetRollingWindowSize{}

type CustomTargetRollingWindowSize struct {
	Value int64 `json:"value"`

	// Fields inherited from TargetRollingWindowSize
}

var _ json.Marshaler = CustomTargetRollingWindowSize{}

func (s CustomTargetRollingWindowSize) MarshalJSON() ([]byte, error) {
	type wrapper CustomTargetRollingWindowSize
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomTargetRollingWindowSize: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomTargetRollingWindowSize: %+v", err)
	}
	decoded["mode"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomTargetRollingWindowSize: %+v", err)
	}

	return encoded, nil
}
