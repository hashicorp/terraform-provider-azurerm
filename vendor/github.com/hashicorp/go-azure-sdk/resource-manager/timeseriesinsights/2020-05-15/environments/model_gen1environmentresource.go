package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentResource = Gen1EnvironmentResource{}

type Gen1EnvironmentResource struct {
	Properties Gen1EnvironmentResourceProperties `json:"properties"`

	// Fields inherited from EnvironmentResource
	Id       *string            `json:"id,omitempty"`
	Location string             `json:"location"`
	Name     *string            `json:"name,omitempty"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

var _ json.Marshaler = Gen1EnvironmentResource{}

func (s Gen1EnvironmentResource) MarshalJSON() ([]byte, error) {
	type wrapper Gen1EnvironmentResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen1EnvironmentResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen1EnvironmentResource: %+v", err)
	}
	decoded["kind"] = "Gen1"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen1EnvironmentResource: %+v", err)
	}

	return encoded, nil
}
