package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileServerLinkedServiceTypeProperties struct {
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	Host                interface{}  `json:"host"`
	Password            SecretBase   `json:"password"`
	UserId              *interface{} `json:"userId,omitempty"`
}

var _ json.Unmarshaler = &FileServerLinkedServiceTypeProperties{}

func (s *FileServerLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		Host                interface{}  `json:"host"`
		UserId              *interface{} `json:"userId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.UserId = decoded.UserId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FileServerLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'FileServerLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
