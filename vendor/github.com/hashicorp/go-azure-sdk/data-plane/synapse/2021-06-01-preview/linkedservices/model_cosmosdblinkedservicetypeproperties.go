package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDbLinkedServiceTypeProperties struct {
	AccountEndpoint     *interface{} `json:"accountEndpoint,omitempty"`
	AccountKey          SecretBase   `json:"accountKey"`
	ConnectionString    *interface{} `json:"connectionString,omitempty"`
	Database            *interface{} `json:"database,omitempty"`
	EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
}

var _ json.Unmarshaler = &CosmosDbLinkedServiceTypeProperties{}

func (s *CosmosDbLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountEndpoint     *interface{} `json:"accountEndpoint,omitempty"`
		ConnectionString    *interface{} `json:"connectionString,omitempty"`
		Database            *interface{} `json:"database,omitempty"`
		EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountEndpoint = decoded.AccountEndpoint
	s.ConnectionString = decoded.ConnectionString
	s.Database = decoded.Database
	s.EncryptedCredential = decoded.EncryptedCredential

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CosmosDbLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accountKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccountKey' for 'CosmosDbLinkedServiceTypeProperties': %+v", err)
		}
		s.AccountKey = impl
	}

	return nil
}
