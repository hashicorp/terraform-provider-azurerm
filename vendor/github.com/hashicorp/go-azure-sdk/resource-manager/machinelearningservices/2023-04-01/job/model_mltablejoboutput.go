package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobOutput = MLTableJobOutput{}

type MLTableJobOutput struct {
	Mode *OutputDeliveryMode `json:"mode,omitempty"`
	Uri  *string             `json:"uri,omitempty"`

	// Fields inherited from JobOutput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = MLTableJobOutput{}

func (s MLTableJobOutput) MarshalJSON() ([]byte, error) {
	type wrapper MLTableJobOutput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MLTableJobOutput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MLTableJobOutput: %+v", err)
	}
	decoded["jobOutputType"] = "mltable"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MLTableJobOutput: %+v", err)
	}

	return encoded, nil
}
