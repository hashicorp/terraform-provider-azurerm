package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = SecretAuthInfo{}

type SecretAuthInfo struct {
	Name       *string        `json:"name,omitempty"`
	SecretInfo SecretInfoBase `json:"secretInfo"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s SecretAuthInfo) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = SecretAuthInfo{}

func (s SecretAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper SecretAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SecretAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretAuthInfo: %+v", err)
	}

	decoded["authType"] = "secret"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SecretAuthInfo: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &SecretAuthInfo{}

func (s *SecretAuthInfo) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Name     *string   `json:"name,omitempty"`
		AuthMode *AuthMode `json:"authMode,omitempty"`
		AuthType AuthType  `json:"authType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Name = decoded.Name
	s.AuthMode = decoded.AuthMode
	s.AuthType = decoded.AuthType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SecretAuthInfo into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["secretInfo"]; ok {
		impl, err := UnmarshalSecretInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SecretInfo' for 'SecretAuthInfo': %+v", err)
		}
		s.SecretInfo = impl
	}

	return nil
}
