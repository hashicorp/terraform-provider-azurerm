package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapEccLinkedServiceTypeProperties struct {
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	Password            SecretBase   `json:"password"`
	Url                 interface{}  `json:"url"`
	Username            *interface{} `json:"username,omitempty"`
}

var _ json.Unmarshaler = &SapEccLinkedServiceTypeProperties{}

func (s *SapEccLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		Url                 interface{}  `json:"url"`
		Username            *interface{} `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Url = decoded.Url
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SapEccLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SapEccLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
