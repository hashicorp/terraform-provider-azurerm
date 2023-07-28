package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetRollingWindowSize = AutoTargetRollingWindowSize{}

type AutoTargetRollingWindowSize struct {

	// Fields inherited from TargetRollingWindowSize
}

var _ json.Marshaler = AutoTargetRollingWindowSize{}

func (s AutoTargetRollingWindowSize) MarshalJSON() ([]byte, error) {
	type wrapper AutoTargetRollingWindowSize
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoTargetRollingWindowSize: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoTargetRollingWindowSize: %+v", err)
	}
	decoded["mode"] = "Auto"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoTargetRollingWindowSize: %+v", err)
	}

	return encoded, nil
}
