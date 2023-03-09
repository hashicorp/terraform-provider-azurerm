package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentCreateOrUpdateParameters = Gen2EnvironmentCreateOrUpdateParameters{}

type Gen2EnvironmentCreateOrUpdateParameters struct {
	Properties Gen2EnvironmentCreationProperties `json:"properties"`

	// Fields inherited from EnvironmentCreateOrUpdateParameters
	Location string             `json:"location"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = Gen2EnvironmentCreateOrUpdateParameters{}

func (s Gen2EnvironmentCreateOrUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper Gen2EnvironmentCreateOrUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen2EnvironmentCreateOrUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen2EnvironmentCreateOrUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Gen2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen2EnvironmentCreateOrUpdateParameters: %+v", err)
	}

	return encoded, nil
}
