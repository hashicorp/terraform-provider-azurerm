package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FileShareConfiguration = MountFileShareConfiguration{}

type MountFileShareConfiguration struct {
	Id                string `json:"id"`
	PrivateEndpointId string `json:"privateEndpointId"`

	// Fields inherited from FileShareConfiguration

	ConfigurationType ConfigurationType `json:"configurationType"`
}

func (s MountFileShareConfiguration) FileShareConfiguration() BaseFileShareConfigurationImpl {
	return BaseFileShareConfigurationImpl{
		ConfigurationType: s.ConfigurationType,
	}
}

var _ json.Marshaler = MountFileShareConfiguration{}

func (s MountFileShareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper MountFileShareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MountFileShareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MountFileShareConfiguration: %+v", err)
	}

	decoded["configurationType"] = "Mount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MountFileShareConfiguration: %+v", err)
	}

	return encoded, nil
}
