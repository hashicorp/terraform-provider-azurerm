package environments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnvironmentCreateOrUpdateParameters = Gen1EnvironmentCreateOrUpdateParameters{}

type Gen1EnvironmentCreateOrUpdateParameters struct {
	Properties Gen1EnvironmentCreationProperties `json:"properties"`

	// Fields inherited from EnvironmentCreateOrUpdateParameters

	Kind     EnvironmentKind    `json:"kind"`
	Location string             `json:"location"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

func (s Gen1EnvironmentCreateOrUpdateParameters) EnvironmentCreateOrUpdateParameters() BaseEnvironmentCreateOrUpdateParametersImpl {
	return BaseEnvironmentCreateOrUpdateParametersImpl{
		Kind:     s.Kind,
		Location: s.Location,
		Sku:      s.Sku,
		Tags:     s.Tags,
	}
}

var _ json.Marshaler = Gen1EnvironmentCreateOrUpdateParameters{}

func (s Gen1EnvironmentCreateOrUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper Gen1EnvironmentCreateOrUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Gen1EnvironmentCreateOrUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Gen1EnvironmentCreateOrUpdateParameters: %+v", err)
	}

	decoded["kind"] = "Gen1"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Gen1EnvironmentCreateOrUpdateParameters: %+v", err)
	}

	return encoded, nil
}
