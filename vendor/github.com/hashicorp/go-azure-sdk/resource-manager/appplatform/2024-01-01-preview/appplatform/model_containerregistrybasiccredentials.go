package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContainerRegistryCredentials = ContainerRegistryBasicCredentials{}

type ContainerRegistryBasicCredentials struct {
	Password string `json:"password"`
	Server   string `json:"server"`
	Username string `json:"username"`

	// Fields inherited from ContainerRegistryCredentials

	Type string `json:"type"`
}

func (s ContainerRegistryBasicCredentials) ContainerRegistryCredentials() BaseContainerRegistryCredentialsImpl {
	return BaseContainerRegistryCredentialsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ContainerRegistryBasicCredentials{}

func (s ContainerRegistryBasicCredentials) MarshalJSON() ([]byte, error) {
	type wrapper ContainerRegistryBasicCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContainerRegistryBasicCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContainerRegistryBasicCredentials: %+v", err)
	}

	decoded["type"] = "BasicAuth"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContainerRegistryBasicCredentials: %+v", err)
	}

	return encoded, nil
}
