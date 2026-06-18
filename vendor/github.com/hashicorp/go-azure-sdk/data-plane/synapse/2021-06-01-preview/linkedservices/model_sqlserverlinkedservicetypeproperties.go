package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlServerLinkedServiceTypeProperties struct {
	ConnectionString    interface{}  `json:"connectionString"`
	EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
	Password            SecretBase   `json:"password"`
	UserName            *interface{} `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SqlServerLinkedServiceTypeProperties{}

func (s *SqlServerLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ConnectionString    interface{}  `json:"connectionString"`
		EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
		UserName            *interface{} `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SqlServerLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SqlServerLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
