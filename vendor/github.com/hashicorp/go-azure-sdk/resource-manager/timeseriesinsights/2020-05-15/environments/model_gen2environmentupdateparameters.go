package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentUpdateParameters = Gen2EnvironmentUpdateParameters{}

type Gen2EnvironmentUpdateParameters struct {
	Properties *Gen2EnvironmentMutableProperties `json:"properties,omitempty"`

	// Fields inherited from EnvironmentUpdateParameters
	Tags *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = Gen2EnvironmentUpdateParameters{}

func (s Gen2EnvironmentUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper Gen2EnvironmentUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen2EnvironmentUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen2EnvironmentUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Gen2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen2EnvironmentUpdateParameters: %+v", err)
	}

	return encoded, nil
}
