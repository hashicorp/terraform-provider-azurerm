package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetLags = AutoTargetLags{}

type AutoTargetLags struct {

	// Fields inherited from TargetLags
}

var _ json.Marshaler = AutoTargetLags{}

func (s AutoTargetLags) MarshalJSON() ([]byte, error) {
	type wrapper AutoTargetLags
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoTargetLags: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoTargetLags: %+v", err)
	}
	decoded["mode"] = "Auto"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoTargetLags: %+v", err)
	}

	return encoded, nil
}
