package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmazonRdsForLinkedServiceTypeProperties struct {
	ConnectionString    interface{} `json:"connectionString"`
	EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
	Password            SecretBase  `json:"password"`
}

var _ json.Unmarshaler = &AmazonRdsForLinkedServiceTypeProperties{}

func (s *AmazonRdsForLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ConnectionString    interface{} `json:"connectionString"`
		EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AmazonRdsForLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'AmazonRdsForLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
