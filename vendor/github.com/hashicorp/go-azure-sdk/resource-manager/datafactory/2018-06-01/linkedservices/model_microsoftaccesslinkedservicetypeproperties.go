package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MicrosoftAccessLinkedServiceTypeProperties struct {
	AuthenticationType  *interface{} `json:"authenticationType,omitempty"`
	ConnectionString    interface{}  `json:"connectionString"`
	Credential          SecretBase   `json:"credential"`
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	Password            SecretBase   `json:"password"`
	UserName            *interface{} `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &MicrosoftAccessLinkedServiceTypeProperties{}

func (s *MicrosoftAccessLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *interface{} `json:"authenticationType,omitempty"`
		ConnectionString    interface{}  `json:"connectionString"`
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		UserName            *interface{} `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MicrosoftAccessLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credential' for 'MicrosoftAccessLinkedServiceTypeProperties': %+v", err)
		}
		s.Credential = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'MicrosoftAccessLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
