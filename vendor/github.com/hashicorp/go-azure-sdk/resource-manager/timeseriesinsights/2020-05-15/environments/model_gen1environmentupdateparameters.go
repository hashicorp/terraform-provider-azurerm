package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentUpdateParameters = Gen1EnvironmentUpdateParameters{}

type Gen1EnvironmentUpdateParameters struct {
	Properties *Gen1EnvironmentMutableProperties `json:"properties,omitempty"`
	Sku        *Sku                              `json:"sku,omitempty"`

	// Fields inherited from EnvironmentUpdateParameters
	Tags *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = Gen1EnvironmentUpdateParameters{}

func (s Gen1EnvironmentUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper Gen1EnvironmentUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen1EnvironmentUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen1EnvironmentUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Gen1"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen1EnvironmentUpdateParameters: %+v", err)
	}

	return encoded, nil
}
