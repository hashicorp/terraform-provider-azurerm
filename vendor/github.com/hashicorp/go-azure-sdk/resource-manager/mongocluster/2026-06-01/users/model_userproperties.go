package users

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserProperties struct {
	IdentityProvider  IdentityProvider   `json:"identityProvider"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Roles             *[]DatabaseRole    `json:"roles,omitempty"`
}

var _ json.Unmarshaler = &UserProperties{}

func (s *UserProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
		Roles             *[]DatabaseRole    `json:"roles,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ProvisioningState = decoded.ProvisioningState
	s.Roles = decoded.Roles

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UserProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["identityProvider"]; ok {
		impl, err := UnmarshalIdentityProviderImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'IdentityProvider' for 'UserProperties': %+v", err)
		}
		s.IdentityProvider = impl
	}

	return nil
}
