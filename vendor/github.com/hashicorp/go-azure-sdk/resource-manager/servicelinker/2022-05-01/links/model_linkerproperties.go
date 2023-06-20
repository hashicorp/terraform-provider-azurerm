package links

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkerProperties struct {
	AuthInfo          AuthInfoBase      `json:"authInfo"`
	ClientType        *ClientType       `json:"clientType,omitempty"`
	ProvisioningState *string           `json:"provisioningState,omitempty"`
	Scope             *string           `json:"scope,omitempty"`
	SecretStore       *SecretStore      `json:"secretStore,omitempty"`
	TargetService     TargetServiceBase `json:"targetService"`
	VNetSolution      *VNetSolution     `json:"vNetSolution,omitempty"`
}

var _ json.Unmarshaler = &LinkerProperties{}

func (s *LinkerProperties) UnmarshalJSON(bytes []byte) error {
	type alias LinkerProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into LinkerProperties: %+v", err)
	}

	s.ClientType = decoded.ClientType
	s.ProvisioningState = decoded.ProvisioningState
	s.Scope = decoded.Scope
	s.SecretStore = decoded.SecretStore
	s.VNetSolution = decoded.VNetSolution

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling LinkerProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authInfo"]; ok {
		impl, err := unmarshalAuthInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthInfo' for 'LinkerProperties': %+v", err)
		}
		s.AuthInfo = impl
	}

	if v, ok := temp["targetService"]; ok {
		impl, err := unmarshalTargetServiceBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TargetService' for 'LinkerProperties': %+v", err)
		}
		s.TargetService = impl
	}
	return nil
}
