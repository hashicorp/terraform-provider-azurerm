package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerRegistryCredentials interface {
	ContainerRegistryCredentials() BaseContainerRegistryCredentialsImpl
}

var _ ContainerRegistryCredentials = BaseContainerRegistryCredentialsImpl{}

type BaseContainerRegistryCredentialsImpl struct {
	Type string `json:"type"`
}

func (s BaseContainerRegistryCredentialsImpl) ContainerRegistryCredentials() BaseContainerRegistryCredentialsImpl {
	return s
}

var _ ContainerRegistryCredentials = RawContainerRegistryCredentialsImpl{}

// RawContainerRegistryCredentialsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawContainerRegistryCredentialsImpl struct {
	containerRegistryCredentials BaseContainerRegistryCredentialsImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawContainerRegistryCredentialsImpl) ContainerRegistryCredentials() BaseContainerRegistryCredentialsImpl {
	return s.containerRegistryCredentials
}

func UnmarshalContainerRegistryCredentialsImplementation(input []byte) (ContainerRegistryCredentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ContainerRegistryCredentials into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BasicAuth") {
		var out ContainerRegistryBasicCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContainerRegistryBasicCredentials: %+v", err)
		}
		return out, nil
	}

	var parent BaseContainerRegistryCredentialsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseContainerRegistryCredentialsImpl: %+v", err)
	}

	return RawContainerRegistryCredentialsImpl{
		containerRegistryCredentials: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
