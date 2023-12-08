package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerRegistryProperties struct {
	Credentials       ContainerRegistryCredentials        `json:"credentials"`
	ProvisioningState *ContainerRegistryProvisioningState `json:"provisioningState,omitempty"`
}

var _ json.Unmarshaler = &ContainerRegistryProperties{}

func (s *ContainerRegistryProperties) UnmarshalJSON(bytes []byte) error {
	type alias ContainerRegistryProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ContainerRegistryProperties: %+v", err)
	}

	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ContainerRegistryProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := unmarshalContainerRegistryCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'ContainerRegistryProperties': %+v", err)
		}
		s.Credentials = impl
	}
	return nil
}
