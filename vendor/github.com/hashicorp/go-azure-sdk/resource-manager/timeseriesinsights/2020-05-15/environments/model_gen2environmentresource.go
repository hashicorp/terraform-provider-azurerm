package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentResource = Gen2EnvironmentResource{}

type Gen2EnvironmentResource struct {
	Properties Gen2EnvironmentResourceProperties `json:"properties"`

	// Fields inherited from EnvironmentResource
	Id       *string            `json:"id,omitempty"`
	Location string             `json:"location"`
	Name     *string            `json:"name,omitempty"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

var _ json.Marshaler = Gen2EnvironmentResource{}

func (s Gen2EnvironmentResource) MarshalJSON() ([]byte, error) {
	type wrapper Gen2EnvironmentResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen2EnvironmentResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen2EnvironmentResource: %+v", err)
	}
	decoded["kind"] = "Gen2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen2EnvironmentResource: %+v", err)
	}

	return encoded, nil
}
