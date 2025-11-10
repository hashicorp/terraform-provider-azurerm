package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapBWLinkedServiceTypeProperties struct {
	ClientId            interface{}  `json:"clientId"`
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	Password            SecretBase   `json:"password"`
	Server              interface{}  `json:"server"`
	SystemNumber        interface{}  `json:"systemNumber"`
	UserName            *interface{} `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SapBWLinkedServiceTypeProperties{}

func (s *SapBWLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId            interface{}  `json:"clientId"`
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		Server              interface{}  `json:"server"`
		SystemNumber        interface{}  `json:"systemNumber"`
		UserName            *interface{} `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Server = decoded.Server
	s.SystemNumber = decoded.SystemNumber
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SapBWLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SapBWLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
